# OpenShift CI Inspector

This project aims to grab the logs, artifacts, etc from [Prow](https://prow.ci.openshift.org/) and index them locally for fast issue inspection.

## Running this project

You can launch the MySQL database and Minio (S3) server using [docker-compose](https://docs.docker.com/compose/install/). (This should work with Podman too.)


Once you have launched the database server you can then apply the [database schema](schema/jobs.sql) using the following command:

```bash
cat schema/jobs.sql | mysql -h 127.0.0.1 -u inspector -pinspector inspector
```

Finally, you can run [scrape_prow.go](scrape_prow.go) using the following command:

```bash
export MYSQL_HOST=127.0.0.1
export MYSQL_USER=inspector
export MYSQL_PASSWORD=inspector
export MYSQL_DB=inspector
export S3_ACCESS_KEY_ID=inspector
export S3_SECRET_ACCESS_KEY=inspector
export S3_ENDPOINT=http://127.0.0.1:9000
export S3_REGION=us-east-1
export S3_PATH_STYLE_ACCESS=1
go run scrape_prow.go
```