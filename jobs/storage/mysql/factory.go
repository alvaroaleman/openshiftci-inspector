package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/janoszen/openshiftci_inspector/common/mysql"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

const (
	// language=MySQL
	createJobsTableSQL = `
CREATE TABLE IF NOT EXISTS jobs
(
    id              VARCHAR(255) PRIMARY KEY COMMENT 'metadata.uid',
    job             VARCHAR(255) COMMENT 'spec.job',
    status          ENUM ('success', 'failure', 'pending', 'aborted', 'error', '') COMMENT 'status.state',
    start_time      DATETIME DEFAULT NULL COMMENT 'status.startTime',
    pending_time    DATETIME DEFAULT NULL COMMENT 'status.pendingTime',
    completion_time DATETIME DEFAULT NULL COMMENT 'status.completionTime',
    url             VARCHAR(255) COMMENT 'status.url',

    job_name_safe   VARCHAR(255),

    git_org         VARCHAR(255),
    git_repo        VARCHAR(255),
    git_repo_link   VARCHAR(255),
    git_base_ref    VARCHAR(255),
    git_base_sha    VARCHAR(255),
    git_base_link   VARCHAR(255),

    artifacts_url VARCHAR(255),

    INDEX i_status (status)
)
`
	// language=MySQL
	createJobsPullsSQL = `
CREATE TABLE IF NOT EXISTS job_pulls
(
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id      VARCHAR(255),
    number      INT,
    author      VARCHAR(255),
    sha         VARCHAR(255),
    pull_link   VARCHAR(255),
    commit_link VARCHAR(255),
    author_link VARCHAR(255),

    INDEX i_job_id (job_id),
    CONSTRAINT fk_job_pulls_job_id
        FOREIGN KEY (job_id)
            REFERENCES jobs (id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT,
    UNIQUE u_pulls (job_id, number, sha)
)
`
)

// NewMySQLJobsStorage creates a MySQL storage for jobs.
func NewMySQLJobsStorage(config mysql.Config) (storage.CompoundJobsStorage, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	db, err := sql.Open(
		"mysql",
		config.ConnectString(),
	)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(
		`CREATE DATABASE IF NOT EXISTS ` + config.Database,
	); err != nil {
		return nil, fmt.Errorf("failed to create database (%w)", err)
	}
	if _, err := db.Exec(createJobsTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create jobs table (%w)", err)
	}
	if _, err := db.Exec(createJobsPullsSQL); err != nil {
		return nil, fmt.Errorf("failed to create job_pulls table (%w)", err)
	}

	return &mysqlJobsStorage{
		db: db,
	}, nil
}
