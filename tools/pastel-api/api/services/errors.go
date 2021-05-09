package services

// API errors.
var (
	ErrNotFoundMethod = NewError(-32601, "method not found")
	ErrNotMasterNode  = NewError(-32603, "this is not a masternode")
)

// Error represents an API error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

// NewError returns a new Error instance.
func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}
