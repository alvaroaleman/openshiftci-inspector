package jobs

import (
	"time"
)

type Job struct {
	ID             string
	Job            string
	JobSafeName    string
	Status         string
	StartTime      *time.Time
	PendingTime    *time.Time
	CompletionTime *time.Time
	URL            string
	GitOrg         string
	GitRepo        string
	GitRepoLink    string
	GitBaseRef     string
	GitBaseSHA     string
	GitBaseLink    string

	Pulls []Pull
}
