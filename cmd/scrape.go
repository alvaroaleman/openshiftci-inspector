package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpDownloader "github.com/janoszen/openshiftci-inspector/asset/downloader/http"
	"github.com/janoszen/openshiftci-inspector/asset/index"
	mysqlAssetIndex "github.com/janoszen/openshiftci-inspector/asset/indexstorage/mysql"
	"github.com/janoszen/openshiftci-inspector/asset/storage/s3"
	commonHTTP "github.com/janoszen/openshiftci-inspector/common/http"
	mysqlCommon "github.com/janoszen/openshiftci-inspector/common/mysql"
	"github.com/janoszen/openshiftci-inspector/jobs/asseturl/caching"
	assetURLHTTP "github.com/janoszen/openshiftci-inspector/jobs/asseturl/http"
	scrapeHTTP "github.com/janoszen/openshiftci-inspector/jobs/scrape/http"
	"github.com/janoszen/openshiftci-inspector/jobs/storage/mysql"
	"github.com/janoszen/openshiftci-inspector/scraper"
)

func main() {
	logger := log.New(os.Stdout, "", 0)

	httpScraper, err := scrapeHTTP.NewHTTPScraper(
		commonHTTP.Config{
			BaseURL: "https://prow.ci.openshift.org/",
			CACert:  os.Getenv("CA_CERT"),
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	mysqlConfig := mysqlCommon.Config{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     3306,
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DB"),
	}

	jobStorage, err := mysql.NewMySQLJobsStorage(mysqlConfig)
	if err != nil {
		panic(err)
	}

	httpAssetURLFetcher, err := assetURLHTTP.NewHTTPAssetURLFetcher(
		commonHTTP.Config{
			BaseURL: "https://prow.ci.openshift.org/",
			CACert:  os.Getenv("CA_CERT"),
		},
	)
	if err != nil {
		panic(err)
	}
	jobsAssetStorage, err := mysql.NewMySQLJobsStorage(mysqlConfig)
	if err != nil {
		panic(err)
	}
	assetURLFetcher := caching.New(httpAssetURLFetcher, jobsAssetStorage, logger)

	assetIndexStorage, err := mysqlAssetIndex.NewMySQLAssetIndex(mysqlConfig, logger)
	if err != nil {
		panic(err)
	}

	assets := map[string][]string{
		"e2e-ovirt": {
			"artifacts/e2e-ovirt/metrics/prometheus.tar",
		},
	}

	assetIndex := index.New(assetIndexStorage, logger, assets)

	pathStyleAccess := os.Getenv("AWS_S3_PATH_STYLE_ACCESS") != ""
	assetStorage, err := s3.NewS3AssetStorage(
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

	assetDownloader, err := httpDownloader.New(
		commonHTTP.Config{
			BaseURL: "https://prow.ci.openshift.org/",
			CACert:  os.Getenv("CA_CERT"),
		},
		assetStorage,
		assetIndexStorage,
		logger,
	)
	if err != nil {
		panic(err)
	}

	program := scraper.New(
		logger,
		httpScraper,
		jobStorage,
		assetURLFetcher,
		assetIndex,
		assetDownloader,
	)

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(
		shutdownSignals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	go func() {
		<-shutdownSignals
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		program.Shutdown(ctx)
	}()

	program.Run()
}
