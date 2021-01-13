package jobs

// Pull is a pull request for a Job.
type Pull struct {
	// Number is the pull request number
	Number int `json:"number"`
	// Author is the GitHub username of the author of the PR.
	Author string `json:"author"`
	// SHA is the SHA of the pull request.
	SHA string `json:"sha"`
	// PullLink is a HTTP URL of the pull request.
	PullLink string `json:"pullLink"`
	// CommitLink is a HTTP URL of the commit.
	CommitLink string `json:"commitLink"`
	// AuthorLink is a HTTP URL of the author's profile.
	AuthorLink string `json:"authorLink"`
}
