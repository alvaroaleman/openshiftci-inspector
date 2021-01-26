package http

import (
	"context"
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
	exit    chan struct{}
	done    chan struct{}
	logger  *log.Logger
}

func (d *assetDownloader) Download(assets <-chan asset.AssetWithJob) {
	go func() {
		for {
			select {
			case <-d.exit:
				return
			case a, ok := <-assets:
				if !ok {
					return
				}
				d.logger.Printf("Downloading asset %s...\n", a.Job.AssetURL+a.AssetName)
				response, err := d.client.Get(a.Job.AssetURL + a.AssetName)
				if err != nil {
					d.logger.Printf("Failed to download URL %s (%v).", a.Job.AssetURL+a.AssetName, err)
					continue
				}
				if response.StatusCode != 200 {
					d.logger.Printf("Failed to download URL %s (non-200 status code: %d).", a.Job.AssetURL+a.AssetName, response.StatusCode)
					continue
				}
				if strings.Contains(response.Header.Get("Content-Type"), "text/html") {
					d.logger.Printf("Failed to download URL %s (non-binary content type: %s).", a.Job.AssetURL+a.AssetName, response.Header.Get("Content-Type"))
					continue
				}
				data, err := ioutil.ReadAll(response.Body)
				_ = response.Body.Close()
				if err != nil {
					d.logger.Printf("Failed to download URL %s (%v).", a.Job.AssetURL+a.AssetName, err)
					continue
				}
				err = d.storage.Store(a.Asset, "application/octet-stream", data)
				if err != nil {
					d.logger.Printf("Failed to store asset %s for job %s (%v).", a.AssetName, a.JobID, err)
					continue
				}
				err = d.index.AddAsset(a.JobID, a.AssetName)
				if err != nil {
					d.logger.Printf("Failed to index asset %s for job %s (%v).", a.AssetName, a.JobID, err)
					continue
				}
			}
		}
	}()
}

func (d *assetDownloader) Shutdown(ctx context.Context) {
	select {
	case <-d.done:
		return
	case <-ctx.Done():
		close(d.exit)
	}
	<-d.done
}
