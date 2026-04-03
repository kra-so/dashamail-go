package dashamail

import "fmt"

// APIError is returned when the DashaMail API responds with a non-zero error code.
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("dashamail: API error %d: %s", e.Code, e.Message)
}
