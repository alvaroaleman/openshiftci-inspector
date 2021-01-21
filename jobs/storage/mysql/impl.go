package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

type mysqlJobsStorage struct {
	db *sql.DB
}

func (m *mysqlJobsStorage) UpdateAssetURL(job jobs.Job, assetURL string) error {
	_, err := m.db.Exec(
		// language=MySQL
		`UPDATE jobs SET artifacts_url=? WHERE id=?`,
		assetURL,
		job.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlJobsStorage) GetAssetURLForJob(job jobs.Job) (assetURL string, err error) {
	result, err := m.db.Query(
		// language=MySQL
		`SELECT artifacts_url FROM jobs WHERE id=?`,
		job.ID,
	)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = result.Close()
	}()
	if !result.Next() {
		return "", storage.JobHasNoAssetURL
	}
	var url *string
	if err := result.Scan(&url); err != nil {
		return "", err
	}
	if url == nil || *url == "" {
		return "", storage.JobHasNoAssetURL
	}
	return *url, nil
}

func (m *mysqlJobsStorage) UpdateJob(job jobs.Job) (err error) {
	if err := m.upsertJob(job); err != nil {
		return err
	}
	parameters := []interface{}{
		job.ID,
	}
	var placeholders []string
	for _, pull := range job.Pulls {
		if err := m.upsertPull(job, pull); err != nil {
			return err
		}
		placeholders = append(placeholders, "?")
		parameters = append(parameters, pull.Number)
	}
	if len(placeholders) > 0 {
		_, err = m.db.Exec(
			`
DELETE FROM job_pulls
WHERE
	job_id=?
	AND
	number NOT IN (`+strings.Join(placeholders, ", ")+`)`,
			parameters...,
		)
		if err != nil {
			return fmt.Errorf("failed to delete unnecessary pulls for job %s (%w)", job.ID, err)
		}
	}
	return nil
}

func (m *mysqlJobsStorage) upsertPull(job jobs.Job, pull jobs.Pull) error {
	_, err := m.db.Exec(
		// language=MySQL
		`
INSERT INTO job_pulls (
  	job_id,
    number,
    author,
    sha,
    pull_link,
    commit_link,
    author_link
) VALUES (
	?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE 
	author = ?,
    sha = ?,
    pull_link = ?,
	commit_link = ?,
	author_link = ?`,
		job.ID,
		pull.Number,
		pull.Author,
		pull.SHA,
		pull.PullLink,
		pull.CommitLink,
		pull.AuthorLink,

		pull.Author,
		pull.SHA,
		pull.PullLink,
		pull.CommitLink,
		pull.AuthorLink,
	)
	if err != nil {
		return fmt.Errorf("failed to insert pull %d for job %s (%w)", pull.Number, job.ID, err)
	}
	return nil
}

func (m *mysqlJobsStorage) upsertJob(job jobs.Job) (err error) {
	_, err = m.db.Exec(
		// language=MySQL
		`
INSERT INTO jobs (
    id,
    job,
    job_name_safe,
    status,
    start_time,
    pending_time,
    completion_time,
    url,
    git_org,
    git_repo,
    git_repo_link,
    git_base_ref,
    git_base_sha,
    git_base_link
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
    status = ?, 
    start_time = ?,
    pending_time = ?,
    completion_time = ?`,
		job.ID,
		job.Job,
		job.JobSafeName,
		job.Status,
		job.StartTime,
		job.PendingTime,
		job.CompletionTime,
		job.URL,
		job.GitOrg,
		job.GitRepo,
		job.GitRepoLink,
		job.GitBaseRef,
		job.GitBaseSHA,
		job.GitBaseLink,

		job.Status,
		job.StartTime,
		job.PendingTime,
		job.CompletionTime,
	)
	if err != nil {
		return fmt.Errorf("failed to insert job %s (%w)", job.ID, err)
	}
	return nil
}

func (m *mysqlJobsStorage) ListJobs() (jobList []jobs.Job, err error) {
	var result *sql.Rows
	result, err = m.db.Query(
		// language=MySQL
		`
SELECT
	id,
    job,
    job_name_safe,
    status,
    start_time,
    pending_time,
    completion_time,
    url,
    git_org,
    git_repo,
    git_repo_link,
    git_base_ref,
    git_base_sha,
    git_base_link
FROM jobs`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs (%w)", err)
	}
	defer func() {
		_ = result.Close()
	}()

	for {
		if !result.Next() {
			break
		}
		job := jobs.Job{}
		var startTime []uint8
		var pendingTime []uint8
		var completionTime []uint8
		err := result.Scan(
			&job.ID,
			&job.Job,
			&job.JobSafeName,
			&job.Status,
			&startTime,
			&pendingTime,
			&completionTime,
			&job.URL,
			&job.GitOrg,
			&job.GitRepo,
			&job.GitRepoLink,
			&job.GitBaseRef,
			&job.GitBaseSHA,
			&job.GitBaseLink,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch job row (%w)", err)
		}
		if job.StartTime, err = m.decodeTime(startTime); err != nil {
			return nil, err
		}
		if job.PendingTime, err = m.decodeTime(pendingTime); err != nil {
			return nil, err
		}
		if job.CompletionTime, err = m.decodeTime(completionTime); err != nil {
			return nil, err
		}

		pulls, err := m.getJobPulls(job.ID)
		if err != nil {
			return nil, err
		}
		job.Pulls = pulls
		jobList = append(jobList, job)
	}
	return jobList, nil
}

func (m *mysqlJobsStorage) getJobPulls(jobID string) ([]jobs.Pull, error) {
	// TODO work around N+1 queries
	pullsResult, err := m.db.Query(
		// language=MySQL
		`
SELECT
	number, author, sha, pull_link, commit_link, author_link
FROM
	job_pulls
WHERE job_id = ?
`,
		jobID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = pullsResult.Close()
	}()
	// We explicitly want this to be an empty slice.
	//goland:noinspection GoPreferNilSlice
	pulls := []jobs.Pull{}
	for {
		if !pullsResult.Next() {
			break
		}

		pull := jobs.Pull{}
		err := pullsResult.Scan(
			&pull.Number,
			&pull.Author,
			&pull.SHA,
			&pull.PullLink,
			&pull.CommitLink,
			&pull.AuthorLink,
		)
		if err != nil {
			return nil, err
		}
		pulls = append(pulls, pull)
	}
	return pulls, nil
}

func (m *mysqlJobsStorage) decodeTime(timeBytes []uint8) (*time.Time, error) {
	if timeBytes != nil {
		t, err := time.Parse("2006-01-02 15:04:05", string(timeBytes))
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	return nil, nil
}

func (m *mysqlJobsStorage) Shutdown(_ context.Context) {
	_ = m.db.Close()
}
