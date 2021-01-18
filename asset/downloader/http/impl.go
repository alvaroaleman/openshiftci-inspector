package http

import (
	"context"
	"io/ioutil"
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
				response, err := d.client.Get(a.Job.AssetURL + a.AssetName)
				if err != nil {
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
