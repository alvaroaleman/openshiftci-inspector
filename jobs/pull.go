package jobs

// Pull is a pull request for a Job.
//
// swagger:model Pull
type Pull struct {
	// Number is the pull request number
	//
	// required: true
	Number int `json:"number"`
	// Author is the GitHub username of the author of the PR.
	//
	// required: true
	Author string `json:"author"`
	// SHA is the SHA of the pull request.
	//
	// required: true
	SHA string `json:"sha"`
	// PullLink is a HTTP URL of the pull request.
	//
	// required: true
	PullLink string `json:"pullLink"`
	// CommitLink is a HTTP URL of the commit.
	//
	// required: true
	CommitLink string `json:"commitLink"`
	// AuthorLink is a HTTP URL of the author's profile.
	//
	// required: true
	AuthorLink string `json:"authorLink"`
}
