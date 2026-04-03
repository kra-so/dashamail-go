package dashamail

import (
	"context"
	"testing"
)

func TestGetTransactionalWebhooks_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{"open":"https://example.com/open"}`))
	c := newTestClient(t, ms)

	data, err := c.GetTransactionalWebhooks(context.Background(), "")
	if err != nil {
		t.Fatalf("GetTransactionalWebhooks: %v", err)
	}
	if data == nil {
		t.Error("data should not be nil")
	}
}

func TestGetTransactionalWebhooks_WithEventName(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{"open":"https://example.com/open"}`))
	c := newTestClient(t, ms)

	_, err := c.GetTransactionalWebhooks(context.Background(), "open")
	if err != nil {
		t.Fatalf("GetTransactionalWebhooks: %v", err)
	}

	if ms.lastBody["event_name"] != "open" {
		t.Errorf("event_name = %v, want %q", ms.lastBody["event_name"], "open")
	}
}

func TestSetTransactionalWebhooks_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	err := c.SetTransactionalWebhooks(context.Background(), &WebhookURLs{
		Open:  "https://example.com/open",
		Click: "https://example.com/click",
		Hard:  "https://example.com/bounce",
	})
	if err != nil {
		t.Fatalf("SetTransactionalWebhooks: %v", err)
	}

	if ms.lastBody["open"] != "https://example.com/open" {
		t.Errorf("open = %v", ms.lastBody["open"])
	}
	if ms.lastBody["click"] != "https://example.com/click" {
		t.Errorf("click = %v", ms.lastBody["click"])
	}
	if ms.lastBody["hard"] != "https://example.com/bounce" {
		t.Errorf("hard = %v", ms.lastBody["hard"])
	}
	// Ensure empty fields are not sent
	if _, exists := ms.lastBody["spam"]; exists {
		t.Error("spam should not be sent when empty")
	}
}

func TestSetTransactionalWebhooks_Nil(t *testing.T) {
	c := New("key")
	err := c.SetTransactionalWebhooks(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil webhooks")
	}
}

func TestDeleteTransactionalWebhooks_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	err := c.DeleteTransactionalWebhooks(context.Background(), "open")
	if err != nil {
		t.Fatalf("DeleteTransactionalWebhooks: %v", err)
	}

	if ms.lastBody["event_name"] != "open" {
		t.Errorf("event_name = %v", ms.lastBody["event_name"])
	}
}

func TestDeleteTransactionalWebhooks_EmptyName(t *testing.T) {
	c := New("key")
	err := c.DeleteTransactionalWebhooks(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty event name")
	}
}

func TestSetTransactionalWebhooks_AllFields(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	err := c.SetTransactionalWebhooks(context.Background(), &WebhookURLs{
		Open:      "https://example.com/open",
		Click:     "https://example.com/click",
		Hard:      "https://example.com/hard",
		Soft:      "https://example.com/soft",
		Spam:      "https://example.com/spam",
		Unsub:     "https://example.com/unsub",
		Subscribe: "https://example.com/subscribe",
		Confirm:   "https://example.com/confirm",
	})
	if err != nil {
		t.Fatalf("SetTransactionalWebhooks: %v", err)
	}

	expected := map[string]string{
		"open":      "https://example.com/open",
		"click":     "https://example.com/click",
		"hard":      "https://example.com/hard",
		"soft":      "https://example.com/soft",
		"spam":      "https://example.com/spam",
		"unsub":     "https://example.com/unsub",
		"subscribe": "https://example.com/subscribe",
		"confirm":   "https://example.com/confirm",
	}
	for key, want := range expected {
		if got, _ := ms.lastBody[key].(string); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}
