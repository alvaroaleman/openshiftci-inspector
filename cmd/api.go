// Package classification Openshift CI Inspector.
//
//     Schemes: http
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
//
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mysqlAssetIndex "github.com/janoszen/openshiftci_inspector/asset/indexstorage/mysql"
	"github.com/janoszen/openshiftci_inspector/asset/storage/s3"
	"github.com/janoszen/openshiftci_inspector/common/api"
	mysqlCommon "github.com/janoszen/openshiftci_inspector/common/mysql"
	jobsAPI "github.com/janoszen/openshiftci_inspector/jobs/api"
	"github.com/janoszen/openshiftci_inspector/jobs/metrics"
	mysqlJobsStorage "github.com/janoszen/openshiftci_inspector/jobs/storage/mysql"
)

func main() {
	logger := log.New(os.Stdout, "", 0)

	mysqlConfig := mysqlCommon.Config{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     3306,
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DB"),
	}

	jobsStorage, err := mysqlJobsStorage.NewMySQLJobsStorage(mysqlConfig)
	if err != nil {
		panic(err)
	}

	pathStyleAccess := os.Getenv("AWS_S3_PATH_STYLE_ACCESS") != ""
	assetStore, err := s3.NewS3AssetStorage(
		s3.S3AssetStorageConfig{
			AccessKey:            os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretKey:            os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Bucket:               os.Getenv("AWS_S3_BUCKET"),
			Region:               os.Getenv("AWS_REGION"),
			Endpoint:             os.Getenv("AWS_S3_ENDPOINT"),
			ForcePathStyleAccess: pathStyleAccess,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	assetIndexStorage, err := mysqlAssetIndex.NewMySQLAssetIndex(mysqlConfig, logger)
	if err != nil {
		panic(err)
	}

	queryBackend := metrics.NewQuery(assetIndexStorage, assetStore)

	handlers := []api.API{
		jobsAPI.NewJobsListAPI(jobsStorage),
		jobsAPI.NewJobsGetAPI(jobsStorage, jobsStorage),
		jobsAPI.NewJobsGetPreviousAPI(jobsStorage),
		jobsAPI.NewJobsGetRelatedAPI(jobsStorage),
		jobsAPI.NewJobsMetricsAPI(jobsStorage, assetStore, queryBackend),
	}

	srv, err := api.NewServer(
		handlers,
		[]api.Encoder{
			api.NewJSONEncoder(),
		},
		[]api.Decoder{
			api.NewPathVarsDecoder(),
			api.NewQueryStringDecoder(),
			api.NewJSONDecoder(),
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(
		shutdownSignals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	if err := srv.Start(); err != nil {
		panic(err)
	}
	<-shutdownSignals
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Stop(ctx)
}
