package asset

import (
	"github.com/janoszen/openshiftci_inspector/jobs"
)

type Asset struct {
	JobID           string `json:"jobID"`
	AssetName       string `json:"assetName"`
	AssetRemotePath string `json:"assetRemotePath"`
}

type AssetWithJob struct {
	Asset `json:",inline"`
	Job   jobs.JobWithAssetURL `json:"job"`
}
