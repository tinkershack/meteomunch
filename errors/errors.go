package errors

const (
	FATAL = "FATAL"
	FAIL  = "FAIL"
)

// Munch error structure for custom errors
//
// Definied for use when it becomes necessary.
type Munch struct {
	Error   error
	Message string
	Code    int
}
