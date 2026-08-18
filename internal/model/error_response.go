package model

// ErrorResponse is the JSON error envelope returned by the API on failure.
type ErrorResponse struct {
	Error string `json:"error"`
}
