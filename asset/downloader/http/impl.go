package http

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

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
					// TODO error handling
					continue
				}
				data, err := ioutil.ReadAll(response.Body)
				_ = response.Body.Close()
				if err != nil {
					// TODO error handling
					continue
				}
				err = d.storage.Store(a.Asset, "application/octet-stream", data)
				if err != nil {
					// TODO error handling
					continue
				}
				err = d.index.AddAsset(a.JobID, a.AssetName)
				if err != nil {
					// TODO error handling
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
