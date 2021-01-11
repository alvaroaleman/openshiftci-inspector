package jobs

type Pull struct {
	Number     int
	Author     string
	SHA        string
	PullLink   string
	CommitLink string
	AuthorLink string
}
