package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janoszen/openshiftci-inspector/common/api"
	mysqlCommon "github.com/janoszen/openshiftci-inspector/common/mysql"
	jobsAPI "github.com/janoszen/openshiftci-inspector/jobs/api"
	mysqlJobsStorage "github.com/janoszen/openshiftci-inspector/jobs/storage/mysql"
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

	handlers := []api.API{
		jobsAPI.NewJobsListAPI(jobsStorage),
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
