package errors

import "net/http"

// struct that describes a REST error
type HTTPError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

// this can be called by validation, or any controller etc
// to standardize errors
func BadRequest(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "Invalid Request",
	}
}

func NotFound(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "Not Found",
	}
}
