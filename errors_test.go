package dashamail

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Code: 4, Message: "Invalid API key"}
	want := "dashamail: API error 4: Invalid API key"
	if e.Error() != want {
		t.Errorf("Error() = %q, want %q", e.Error(), want)
	}
}

func TestAPIError_ErrorsAs(t *testing.T) {
	original := &APIError{Code: 10, Message: "test"}
	wrapped := error(original)

	var target *APIError
	if !errors.As(wrapped, &target) {
		t.Error("errors.As should match *APIError")
	}
	if target.Code != 10 {
		t.Errorf("Code = %d, want 10", target.Code)
	}
}

func TestAPIError_Codes(t *testing.T) {
	tests := []struct {
		code int
		msg  string
	}{
		{0, "OK"},
		{1, "Unknown error"},
		{4, "Invalid API key"},
		{100, "Rate limit exceeded"},
	}

	for _, tt := range tests {
		e := &APIError{Code: tt.code, Message: tt.msg}
		if e.Code != tt.code {
			t.Errorf("Code = %d, want %d", e.Code, tt.code)
		}
		if e.Message != tt.msg {
			t.Errorf("Message = %q, want %q", e.Message, tt.msg)
		}
	}
}
