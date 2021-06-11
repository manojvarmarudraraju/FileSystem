package web

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error is used to add web information to a request error.
// application with web specific context.
type AppError struct {
	Err error
	Status int
	Fields []FieldError
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (a AppError) Error() string {
	return a.Err.Error()
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, statusCode int) error  {
	return &AppError{
		Err:    err,
		Status: statusCode,
	}
}


