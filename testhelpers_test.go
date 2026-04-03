package dashamail

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockServer creates an httptest.Server that captures requests and returns
// configured responses. The caller must call close() when done.
type mockServer struct {
	server *httptest.Server

	// lastMethod is the HTTP method of the last request.
	lastMethod string
	// lastPath is the full URL path+query of the last request.
	lastPath string
	// lastBody is the parsed JSON body of the last request.
	lastBody map[string]any
	// lastHeaders holds the headers of the last request.
	lastHeaders http.Header

	// responseCode is the HTTP status code to return (default 200).
	responseCode int
	// responseBody is the raw JSON response to return.
	responseBody string
}

func newMockServer(t *testing.T) *mockServer {
	t.Helper()
	ms := &mockServer{
		responseCode: http.StatusOK,
	}

	ms.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.lastMethod = r.Method
		ms.lastPath = r.URL.String()
		ms.lastHeaders = r.Header.Clone()

		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err == nil && len(body) > 0 {
				ms.lastBody = make(map[string]any)
				_ = json.Unmarshal(body, &ms.lastBody)
			}
		}

		code := ms.responseCode
		if code == 0 {
			code = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if ms.responseBody != "" {
			_, _ = w.Write([]byte(ms.responseBody))
		}
	}))

	t.Cleanup(func() { ms.server.Close() })
	return ms
}

func (ms *mockServer) url() string {
	return ms.server.URL
}

func (ms *mockServer) setResponse(code int, body string) {
	ms.responseCode = code
	ms.responseBody = body
}

// okResponse builds a standard success response envelope.
func okResponse(data string) string {
	return `{"response":{"msg":{"err_code":0,"text":"OK","type":"message"},"data":` + data + `}}`
}

// errResponse builds an error response envelope.
func errResponse(code int, text string) string {
	return `{"response":{"msg":{"err_code":` + itoa(code) + `,"text":"` + text + `","type":"message"},"data":{}}}`
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

// newTestClient creates a Client pointing to the mock server.
func newTestClient(t *testing.T, ms *mockServer, opts ...Option) *Client {
	t.Helper()
	defaults := []Option{WithEndpoint(ms.url())}
	defaults = append(defaults, opts...)
	return New("test-api-key", defaults...)
}
