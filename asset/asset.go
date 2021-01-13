package asset

import (
	"github.com/janoszen/openshiftci-inspector/jobs"
)

type Asset struct {
	JobID     string `json:"jobID"`
	AssetName string `json:"assetName"`
}

type AssetWithJob struct {
	Asset `json:",inline"`
	Job   jobs.JobWithAssetURL `json:"job"`
}
