package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/go-sql-driver/mysql"
	v1 "k8s.io/api/core/v1"
)

type JobMetadata struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	SelfLink          string    `json:"selfLink"`
	UID               string    `json:"uid"`
	ResourceVersion   string    `json:"resourceVersion"`
	Generation        int       `json:"generation"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type JobSpec struct {
	Type         string     `json:"type"`
	Agent        string     `json:"agent"`
	Cluster      string     `json:"cluster"`
	Namespace    string     `json:"namespace"`
	Job          string     `json:"job"`
	Refs         JobRefs    `json:"refs"`
	Report       bool       `json:"report"`
	Context      string     `json:"context"`
	RerunCommand string     `json:"rerun_command"`
	PodSpec      v1.PodSpec `json:"pod_spec"`
}

type JobRefs struct {
	Org      string    `json:"org"`
	Repo     string    `json:"repo"`
	RepoLink string    `json:"repo_link"`
	BaseRef  string    `json:"base_ref"`
	BaseSha  string    `json:"base_sha"`
	BaseLink string    `json:"base_link"`
	Pulls    []JobPull `json:"pulls"`
}

type JobStatus struct {
	StartTime      *time.Time `json:"startTime,omitempty"`
	PendingTime    *time.Time `json:"pendingTime,omitempty"`
	CompletionTime *time.Time `json:"completionTime,omitempty"`
	State          string     `json:"state"`
	Description    string     `json:"description"`
	URL            string     `json:"url"`
	PodName        string     `json:"pod_name"`
	BuildID        string     `json:"build_id"`
}

type JobPull struct {
	Number     int    `json:"number"`
	Author     string `json:"author"`
	SHA        string `json:"sha"`
	Link       string `json:"link"`
	CommitLink string `json:"commit_link"`
	AuthorLink string `json:"author_link"`
}

type Job struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   JobMetadata `json:"metadata"`
	Spec       JobSpec     `json:"spec"`
	Status     JobStatus   `json:"status"`
}

type Jobs struct {
	Items []Job `json:"items"`
}

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func mustTX(err error, tx *sql.Tx) {
	if err != nil {
		_ = tx.Rollback()
		log.Fatalln(err)
	}
}

func fetchJobs() ([]Job, error) {
	log.Println("Fetching jobs from Prow...")
	url := "https://prow.ci.openshift.org/prowjobs.js?var=allBuilds"
	client := http.Client{}
	data, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	rawJSON := strings.Replace(string(b[:len(b)-1]), "var allBuilds = ", "", 1)

	jobsRaw := &map[string]interface{}{}
	if err := json.Unmarshal([]byte(rawJSON), jobsRaw); err != nil {
		return nil, err
	}

	jobs := Jobs{}
	if err := json.Unmarshal([]byte(rawJSON), &jobs); err != nil {
		return nil, err
	}
	return jobs.Items, err
}

func main() {
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDB := os.Getenv("MYSQL_DB")

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlDB),
	)
	must(err)

	awsConfig := &aws.Config{
		Credentials: credentials.NewCredentials(
			&credentials.StaticProvider{
				Value: credentials.Value{
					AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),

					SessionToken: "",
					ProviderName: "",
				},
			},
		),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		Region:           aws.String(os.Getenv("S3_REGION")),
		S3ForcePathStyle: aws.Bool(os.Getenv("S3_PATH_STYLE_ACCESS") == "1"),
	}
	sess := session.Must(session.NewSession(awsConfig))
	s3connection := s3.New(sess)
	updateJobs(db)
	updateArtifactURLs(db)
	bucket := os.Getenv("S3_BUCKET")

	updateAssets(
		db,
		s3connection,
		bucket,
		"build_log",
		"build-log.txt",
		"text/plain",
		"ovirt",
	)

	updateAssets(
		db,
		s3connection,
		bucket,
		"events",
		"artifacts/e2e-ovirt/events.json",
		"application/json",
		"ovirt",
	)

	updateAssets(
		db,
		s3connection,
		bucket,
		"prometheus",
		"artifacts/e2e-ovirt/metrics/prometheus.tar",
		"application/json",
		"ovirt",
	)
}

func updateAssets(
	db *sql.DB,
	s3Connection *s3.S3,
	bucketName string,
	assetType string,
	sourcePath string,
	mimeType string,
	filter string,
) {
	log.Println("Fetching " + assetType + " assets...")
	var res *sql.Rows
	var err error
	if filter != "" {
		filterLike := "%" + filter + "%"
		res, err = db.Query(
			"SELECT jobs.id, artifacts_url FROM jobs LEFT JOIN assets a on jobs.id = a.job_id AND asset_type=? WHERE a.job_id IS NULL AND (url LIKE ? OR job_name_safe LIKE ? OR job LIKE ?)",
			assetType, filterLike, filterLike, filterLike,
		)
	} else {
		res, err = db.Query(
			"SELECT jobs.id, artifacts_url FROM jobs LEFT JOIN assets a on jobs.id = a.job_id AND asset_type=? WHERE a.job_id IS NULL",
			assetType,
		)
	}
	must(err)

	_, err = s3Connection.GetBucketLocation(
		&s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		},
	)
	if err != nil {
		_, err = s3Connection.CreateBucket(
			&s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
				ACL:    aws.String(s3.BucketCannedACLPublicRead),
			},
		)
		must(err)
	}

	lock := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}
	for {
		id := ""
		url := ""
		res.Next()
		err := res.Scan(&id, &url)
		if err != nil {
			break
		}
		if url == "" {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock <- struct{}{}
			defer func() { <-lock }()
			log.Println("Fetching the " + assetType + " asset for job " + id + "...")

			response, err := http.Get(url + sourcePath)
			if err != nil {
				log.Printf("Failed to fetch asset %s for job %s", assetType, id)
				return
			}
			data, err := ioutil.ReadAll(response.Body)
			must(err)
			must(response.Body.Close())

			if len(data) == 0 {
				return
			}

			key := "/" + id + "/" + assetType
			_, err = s3Connection.PutObject(
				&s3.PutObjectInput{
					ACL:           aws.String(s3.BucketCannedACLPublicRead),
					Body:          bytes.NewReader(data),
					Bucket:        aws.String(bucketName),
					ContentLength: aws.Int64(int64(len(data))),
					ContentType:   aws.String(mimeType),
					Key:           aws.String(key),
				},
			)
			must(err)
			req, _ := s3Connection.GetObjectRequest(
				&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(key),
				},
			)
			must(req.Sign())
			assetURL := req.HTTPRequest.URL.String()

			insertRes, err := db.Query(
				"INSERT INTO assets (job_id, asset_type, asset_key, asset_url) VALUES (?, ?, ?, ?)",
				id, assetType, key, assetURL,
			)
			must(err)
			must(insertRes.Close())
		}()
	}
	must(res.Close())
	wg.Wait()
}

func updateArtifactURLs(db *sql.DB) {
	res, err := db.Query("SELECT id, url FROM jobs WHERE artifacts_url IS NULL")
	must(err)

	artifactsRe := regexp.MustCompile(`<a href="(?P<url>[^"]+)">Artifacts</a>`)

	lock := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}
	for {
		id := ""
		url := ""
		res.Next()
		err := res.Scan(&id, &url)
		if err != nil {
			break
		}
		if url == "" {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock <- struct{}{}
			log.Printf("Fetching artifact URL for job %s...\n", id)
			jobPage, err := http.Get(url)
			must(err)
			body, err := ioutil.ReadAll(jobPage.Body)
			must(err)
			matches := artifactsRe.FindStringSubmatch(string(body))
			if len(matches) > 1 {
				updateRes, err := db.Query(
					"UPDATE jobs SET artifacts_url=? WHERE id=?", matches[1], id,
				)
				must(err)
				must(updateRes.Close())
			}
			<-lock
		}()
	}
	must(res.Close())
	wg.Wait()
}

func updateJobs(db *sql.DB) {
	jobs, err := fetchJobs()
	must(err)

	for _, job := range jobs {
		tx, err := db.BeginTx(context.TODO(), &sql.TxOptions{})
		must(err)
		jobNameSafe := ""
		if len(job.Spec.PodSpec.Containers) > 0 {
			for _, env := range job.Spec.PodSpec.Containers[0].Env {
				if env.Name == "JOB_NAME_SAFE" {
					jobNameSafe = env.Value
				}
			}
		}
		_, err = tx.Exec(
			"INSERT INTO jobs ("+
				"id, job, job_name_safe, status, start_time, pending_time, completion_time, url,"+
				"git_org, git_repo, git_repo_link, git_base_ref, git_base_sha, git_base_link"+
				") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"+
				"ON DUPLICATE KEY UPDATE "+
				"status = ?,"+
				"start_time = ?,"+
				"pending_time = ?,"+
				"completion_time = ?",
			job.Metadata.UID,
			job.Spec.Job,
			jobNameSafe,
			job.Status.State,
			job.Status.StartTime,
			job.Status.PendingTime,
			job.Status.CompletionTime,
			job.Status.URL,
			job.Spec.Refs.Org,
			job.Spec.Refs.Repo,
			job.Spec.Refs.RepoLink,
			job.Spec.Refs.BaseRef,
			job.Spec.Refs.BaseSha,
			job.Spec.Refs.BaseLink,
			job.Status.State,
			job.Status.StartTime,
			job.Status.PendingTime,
			job.Status.CompletionTime,
		)
		must(err)

		//region Pulls
		result, err := tx.Query("SELECT COUNT(*) AS cnt FROM pulls WHERE job_id=?", job.Metadata.UID)
		mustTX(err, tx)
		pulls := 0
		result.Next()
		err = result.Scan(&pulls)
		mustTX(err, tx)
		err = result.Close()
		mustTX(err, tx)

		if pulls == 0 {
			for _, pull := range job.Spec.Refs.Pulls {
				_, err = tx.Query(
					"REPLACE INTO pulls ("+
						"job_id, number, author, sha, pull_link, commit_link, author_link"+
						") VALUES ("+
						"?, ?, ?, ?, ?, ?, ?"+
						")",
					job.Metadata.UID, pull.Number, pull.Author, pull.SHA, pull.Link, pull.CommitLink, pull.AuthorLink,
				)
				mustTX(err, tx)
			}
		}
		//endregion

		if err := tx.Commit(); err != nil {
			log.Fatalln(err)
		}
	}
}
