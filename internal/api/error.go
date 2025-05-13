package api

import (
	"fmt"
)

// APIError represent an error interacting with the API.
type APIError struct {
	Summary string
	Detail  string
}

// Error returns the error message.
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Summary, e.Detail)
}
