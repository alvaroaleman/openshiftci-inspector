# OpenShift CI Inspector

This project aims to grab the logs, artifacts, etc from [Prow](https://prow.ci.openshift.org/) and index them locally for fast issue inspection.

## Running this project

You can launch the MySQL database and Minio (S3) server using [docker-compose](https://docs.docker.com/compose/install/). (This should work with Podman too.)

### Running the scraper 

Once you have a MySQL and Minio up and running you can run [cmd/scrape.go](cmd/scrape) using the following command:

```bash
export MYSQL_HOST=127.0.0.1
export MYSQL_USER=inspector
export MYSQL_PASSWORD=inspector
export MYSQL_DB=inspector
export AWS_ACCESS_KEY_ID=inspector
export AWS_SECRET_ACCESS_KEY=inspector
export AWS_S3_ENDPOINT=http://127.0.0.1:9000
export AWS_REGION=us-east-1
export AWS_S3_PATH_STYLE_ACCESS=1
export AWS_S3_BUCKET=inspector
go run cmd/scrape.go
```

### Building the frontend

In order to access the web frontend you will need Node and NPM installed. You can then build the frontend code by running:

```bash
npm install
npm run build
```

### Running the frontend

You can also run the frontend using the `npm start` command. API calls will be proxied to the API on port 8080 (see below). 

### Running the API

The API and web service can be run similar to the code above:

```bash
export MYSQL_HOST=127.0.0.1
export MYSQL_USER=inspector
export MYSQL_PASSWORD=inspector
export MYSQL_DB=inspector
export AWS_ACCESS_KEY_ID=inspector
export AWS_SECRET_ACCESS_KEY=inspector
export AWS_S3_ENDPOINT=http://127.0.0.1:9000
export AWS_REGION=us-east-1
export AWS_S3_PATH_STYLE_ACCESS=1
export AWS_S3_BUCKET=inspector
go run cmd/api.go
```

This will start the API on `0.0.0.0:8080`. You can test it by going to http://localhost:8080/ .

### Generating the client library

You can use [Go Swagger](https://goswagger.io/) to generate the Swagger definitions for the API:

```
goswagger generate spec -o swagger.json
```

You can then take the swagger file to generate a Typescript frontend client:

```
cd frontend
npm run client
```

This will update the src/api.ts file from the swagger.json file generated above.
