package code

// Code represents Error Code
type Code string

const (
	// OK shows it is ok.
	OK Code = "OK"

	// Forbidden shows forbidden.
	Forbidden Code = "forbidden"

	// NotFound shows not found.
	NotFound Code = "not_found"
	// InvalidArgument shows the argument is invalid.
	InvalidArgument Code = "invalid_argument"

	// Unknown shows the error is unknown.
	Unknown Code = "unknown"
	// Unexpected shows there's unexpected error.
	Unexpected Code = "unexpected"

	// ContentExpired shows expired term.
	ContentExpired Code = "content_expired"
)
