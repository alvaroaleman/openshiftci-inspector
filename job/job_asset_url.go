package job

// JobWithAssetURL is a job with an already-scraped asset URL.
//
// swagger:model JobWithAssetURL
type JobWithAssetURL struct {
	Job `json:",inline"`

	// AssetURL is the base URL for all assets.
	AssetURL string `json:"assetURL,omitEmpty"`
}
