package jobs

type JobWithAssetURL struct {
	Job `json:",inline"`

	AssetURL string `json:"assetURL"`
}
