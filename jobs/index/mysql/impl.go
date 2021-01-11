package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

type mysqlJobsIndex struct {
	db *sql.DB
}

func (m *mysqlJobsIndex) UpdateJob(job jobs.Job) (err error) {
	if err := m.upsertJob(job); err != nil {
		return err
	}
	var pullNumbers []int
	var placeholders []string
	for _, pull := range job.Pulls {
		if err := m.upsertPull(job, pull); err != nil {
			return err
		}
		pullNumbers = append(pullNumbers, pull.Number)
		placeholders = append(placeholders, "?")
	}
	_, err = m.db.Exec(
		`-- language=mysql
DELETE FROM job_pulls
WHERE
	job_id=?
	AND
	number NOT IN (`+strings.Join(placeholders, ", ")+`)`,
		append([]interface{}{job.ID}, pullNumbers)...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete unnecessary pulls for job %s (%w)", job.ID, err)
	}
	return nil
}

func (m *mysqlJobsIndex) upsertPull(job jobs.Job, pull jobs.Pull) error {
	_, err := m.db.Exec(
		`-- language=mysql
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

func (m *mysqlJobsIndex) upsertJob(job jobs.Job) (err error) {
	_, err = m.db.Exec(
		`-- language=mysql
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

func (m *mysqlJobsIndex) ListJobs() (jobList []jobs.Job, err error) {
	var result *sql.Rows
	result, err = m.db.Query(
		`-- language=mysql
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
	for {
		if !result.Next() {
			break
		}
		job := jobs.Job{}
		err := result.Scan(
			&job.ID,
			&job.Job,
			&job.JobSafeName,
			&job.Status,
			&job.StartTime,
			&job.PendingTime,
			&job.CompletionTime,
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
		jobList = append(jobList, job)
	}
	return jobList, nil
}

func (m *mysqlJobsIndex) Shutdown(_ context.Context) {
	_ = m.db.Close()
}
