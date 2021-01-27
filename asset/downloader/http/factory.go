package http

import (
	"log"

	"github.com/janoszen/openshiftci-inspector/asset/downloader"
	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci-inspector/asset/storage"
	"github.com/janoszen/openshiftci-inspector/common/http"
)

func New(
	config http.Config,
	storage storage.AssetStorage,
	index indexstorage.AssetIndex,
	logger *log.Logger,
) (
	downloader.AssetDownloader,
	error,
) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	client, err := config.CreateClient()
	if err != nil {
		return nil, err
	}

	return &assetDownloader{
		client:  client,
		storage: storage,
		index:   index,
		logger:  logger,
	}, nil
}
