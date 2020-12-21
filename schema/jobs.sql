DROP TABLE IF EXISTS refs;
DROP TABLE IF EXISTS pulls;
DROP TABLE IF EXISTS build_logs;
DROP TABLE IF EXISTS jobs;
CREATE TABLE jobs
(
    id              VARCHAR(255) PRIMARY KEY COMMENT 'metadata.uid',
    job             VARCHAR(255) COMMENT 'spec.job',
    status          ENUM ("success", "failure", "pending", "aborted", "error", "") COMMENT 'status.state',
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
);

CREATE TABLE pulls
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
    CONSTRAINT fk_refs_job_id_jobs
        FOREIGN KEY (job_id)
            REFERENCES jobs (id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT,
    UNIQUE u_pulls (job_id, number, sha)
);

CREATE TABLE build_logs (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id      VARCHAR(255),
    log         LONGBLOB,
    UNIQUE i_job_id(job_id),
    CONSTRAINT fk_build_logs_job_id
        FOREIGN KEY fk_build_logs_job_id(job_id)
            REFERENCES jobs (id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT
)
