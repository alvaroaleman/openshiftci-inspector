package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDB := os.Getenv("MYSQL_DB")

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlDB),
	)
	if err != nil {
		log.Fatalln(err)
	}

	jobs, err := FetchJobs()
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, job := range jobs {
		tx, err := db.BeginTx(context.TODO(), &sql.TxOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		if _, err := tx.Exec(
			"REPLACE INTO jobs (" +
				"id, job, status, start_time, pending_time, completion_time, url," +
				"git_org, git_repo, git_repo_link, git_base_ref, git_base_sha, git_base_link" +
				") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			job.Metadata.UID,
			job.Spec.Job,
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
		); err != nil {
			_ = tx.Rollback()
			log.Fatalln(err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalln(err)
		}
	}
}
