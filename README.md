# OpenShift CI Inspector

This project aims to grab the logs, artifacts, etc from [Prow](https://prow.ci.openshift.org/) and index them locally for fast issue inspection.

## Development

This project consists of 3 main parts:

1. The scraper
2. The API server
3. The frontend

All components require a running MySQL database and an S3-compatible object storage to store assets. The easiest way to do this is to run `docker-compose up -d` in the project directory.

Once these are running you should set the following environment variables:

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
```

### Developing the scraper

The scraper is launched from [cmd/scrape.go](cmd/scrape.go) (`go run cmd/scrape.go`) and requires Golang 1.14 or later to work. It can be run in your usual Go development environment.

### Developing the API

The API is a regular Go HTTP server. It uses the new Golang 1.16 [embed feature](https://github.com/golang/go/issues/41191). **You must install Golang 1.16 to develop the API,** but NodeJS is not required to work on the API.

You can run the API by running [cmd/api.go](cmd/api.go). (`go run cmd/api.go`) The web UI and the API will then run on port 8080.

The backend exposes Swagger/OpenAPI API that is used to automatically generate the Typescript library used by the frontend. You can trigger exporting the `swagger.json` file by running `go run buildtool.go export-api`. The Swagger export is facilitated by [Goswagger](https://goswagger.io/).

This only generates the Swagger file, but does not run the Typescript generation. In order to run the Swagger export and Typescript generation in one step (requires NPM) please run `go run buildtool.go api`.

### Developing the frontend

The frontend is built in Typescript and React.js using the [Material-UI framework](https://material-ui.com/). In order to work on it you will need a working [NodeJS and NPM](https://nodejs.org/en/) installation.

The following commands can be used in the [frontend](frontend) folder.

**Note:** Since most people working on the API are not comfortable with having NodeJS installed the repository contains the built frontend. Please update the build folder regularly.


#### To install dependencies:

```
npm install
```

### To start the dev server

```
npm start
```

### To generate the Typescript client files

```
npm run client
```

### To build a production-ready frontend

```
npm run build
```

## Building everything

In order to build everything (frontend and backend) you can run `go generate && go build -o scrape cmd/scrape.go && go build -o api cmd/api.go`.

## Building for production

Currently, this tool has no production build yet. A preliminary [`Dockerfile`](Dockerfile) is provided.