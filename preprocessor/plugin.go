package preprocessor

import (
	"io"

	"github.com/janoszen/openshiftci_inspector/job"
	"github.com/janoszen/openshiftci_inspector/widget"
)

// Plugin is a means by which artifacts can be transformed into data that can be consumed by frontend plugins.
// It is strongly encouraged that this plugin should only rely on artifact. In functional programming terms this should
// be a pure component.
type Plugin interface {
	// GetArtifacts returns the list of assets this plugin needs to construct the required output.
	//
	// The returned list of assets may specify globs to match files that can be present in multiple locations.
	GetArtifacts(job job.Job) []string

	// Preprocess is responsible for creating a frontend-consumable result from zero or more assets. The assets are
	// provided as map of file names and readers the data can be read from. The preprocessor must provide one or more
	// widgets as a response or an error. It is strongly encouraged to provide a PluginError for errors so that the
	// runner can decide if the error can be retried. If a different error is provided it is assumed that the
	// preprocessing can be retried.
	Preprocess(job job.Job, artifacts map[string]io.Reader) (map[string]widget.Widget, error)
}

// PluginError is an error describing an error that happened during the Preprocess method of the Plugin interface.
type PluginError interface {
	// IsPermanent returns true if the preprocessing should not be retried.
	IsPermanent() bool
	// Error returns the string-error.
	Error() string
}
