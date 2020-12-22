# OpenShift CI Inspector

This project aims to grab the logs, artifacts, etc from [Prow](https://prow.ci.openshift.org/) and index them locally for fast issue inspection.

## Running this project

You can launch the MySQL database and Minio (S3) server using [docker-compose](https://docs.docker.com/compose/install/). (This should work with Podman too.)

Once you have launched the database server you can run the scraping process by running `go run scrape_prow.go`.