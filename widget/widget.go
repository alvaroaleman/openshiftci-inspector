package widget

// Widget is a server-side JSON serializable structure that can be used to render a front-end component.
type Widget interface {
	// Type returns the ID of the widget this plugin is producing output for. The implementation must also be JSON
	// serializable for the frontend to consume.
	Type() string

	// Validate validates the data provided for completeness. This method is called before sending the data to the
	// UI.
	Validate() error

	// JSON must provide a JSON-serializable data structure consumable to a front-end component.
	JSON() interface{}
}
