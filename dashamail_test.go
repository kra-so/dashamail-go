package dashamail

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	c := New("my-key")

	if c.apiKey != "my-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "my-key")
	}
	if c.endpoint != DefaultEndpoint {
		t.Errorf("endpoint = %q, want %q", c.endpoint, DefaultEndpoint)
	}
	if !c.noTrackOpens {
		t.Error("noTrackOpens should default to true")
	}
	if !c.noTrackClicks {
		t.Error("noTrackClicks should default to true")
	}
	if c.ignoreDeliveryPolicy {
		t.Error("ignoreDeliveryPolicy should default to false")
	}
	if c.fromEmail != "" {
		t.Errorf("fromEmail = %q, want empty", c.fromEmail)
	}
	if c.fromName != "" {
		t.Errorf("fromName = %q, want empty", c.fromName)
	}
	if c.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestNew_WithOptions(t *testing.T) {
	custom := &http.Client{Timeout: 60 * time.Second}
	c := New("key",
		WithEndpoint("https://custom.api"),
		WithFromEmail("from@test.com"),
		WithFromName("Test Sender"),
		WithNoTrackOpens(false),
		WithNoTrackClicks(false),
		WithIgnoreDeliveryPolicy(true),
		WithHTTPClient(custom),
		WithDebug(true),
	)

	if c.endpoint != "https://custom.api" {
		t.Errorf("endpoint = %q, want %q", c.endpoint, "https://custom.api")
	}
	if c.fromEmail != "from@test.com" {
		t.Errorf("fromEmail = %q", c.fromEmail)
	}
	if c.fromName != "Test Sender" {
		t.Errorf("fromName = %q", c.fromName)
	}
	if c.noTrackOpens {
		t.Error("noTrackOpens should be false")
	}
	if c.noTrackClicks {
		t.Error("noTrackClicks should be false")
	}
	if !c.ignoreDeliveryPolicy {
		t.Error("ignoreDeliveryPolicy should be true")
	}
	if c.httpClient != custom {
		t.Error("httpClient should be the custom client")
	}
	if !c.debug {
		t.Error("debug should be true")
	}
}

func TestDo_SetsHeaders(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	_, _ = c.do(context.Background(), "test.method", http.MethodPost, map[string]any{"foo": "bar"})

	if ms.lastHeaders.Get("User-Agent") != userAgent {
		t.Errorf("User-Agent = %q, want %q", ms.lastHeaders.Get("User-Agent"), userAgent)
	}
	if ms.lastHeaders.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q", ms.lastHeaders.Get("Content-Type"))
	}
}

func TestDo_SetsAPIKeyInBody(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	_, _ = c.do(context.Background(), "test.method", http.MethodPost, map[string]any{"foo": "bar"})

	if ms.lastBody["api_key"] != "test-api-key" {
		t.Errorf("api_key = %v, want %q", ms.lastBody["api_key"], "test-api-key")
	}
}

func TestDo_SetsMethodQueryParam(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	_, _ = c.do(context.Background(), "transactional.send", http.MethodPost, map[string]any{})

	// Check that ?method=transactional.send is in the URL
	if ms.lastPath == "" {
		t.Fatal("no request was made")
	}
	want := "method=transactional.send"
	if !containsSubstring(ms.lastPath, want) {
		t.Errorf("URL %q does not contain %q", ms.lastPath, want)
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
