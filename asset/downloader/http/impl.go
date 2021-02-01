package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci-inspector/asset/storage"
)

type assetDownloader struct {
	client  *http.Client
	storage storage.AssetStorage
	index   indexstorage.AssetIndex
	logger  *log.Logger
}

func (d *assetDownloader) Download(a asset.AssetWithJob) error {
	response, err := d.client.Get(a.Job.AssetURL + a.AssetRemotePath)
	if err != nil {
		return fmt.Errorf("failed to download URL %s (%w)", a.Job.AssetURL+a.AssetName, err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf(
			"failed to download URL %s (non-200 status code: %d)",
			a.Job.AssetURL+a.AssetName,
			response.StatusCode,
		)
	}
	if strings.Contains(response.Header.Get("Content-Type"), "text/html") {
		return fmt.Errorf(
			"failed to download URL %s (non-binary content type: %s)",
			a.Job.AssetURL+a.AssetName,
			response.Header.Get("Content-Type"),
		)
	}
	data, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to download URL %s (%w)", a.Job.AssetURL+a.AssetName, err)
	}
	err = d.storage.Store(a.Asset, "application/octet-stream", data)
	if err != nil {
		return fmt.Errorf("failed to store asset %s for job %s (%w)", a.AssetName, a.JobID, err)
	}
	err = d.index.AddAsset(a.JobID, a.AssetName)
	if err != nil {
		return fmt.Errorf("failed to index asset %s for job %s (%w)", a.AssetName, a.JobID, err)
	}
	return nil
}
