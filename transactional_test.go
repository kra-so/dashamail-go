package dashamail

import (
	"context"
	"testing"
)

func TestSend_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{"transaction_id":"abc123"}`))
	c := newTestClient(t, ms, WithFromEmail("from@test.com"))

	resp, err := c.Send(context.Background(), &Message{
		To:      "user@test.com",
		Subject: "Test",
		HTML:    "<p>Hello</p>",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	if resp.TransactionID != "abc123" {
		t.Errorf("TransactionID = %q, want %q", resp.TransactionID, "abc123")
	}

	// Verify the request was POST
	if ms.lastMethod != "POST" {
		t.Errorf("method = %q, want POST", ms.lastMethod)
	}

	// Verify body fields
	if ms.lastBody["to"] != "user@test.com" {
		t.Errorf("to = %v", ms.lastBody["to"])
	}
	if ms.lastBody["subject"] != "Test" {
		t.Errorf("subject = %v", ms.lastBody["subject"])
	}
	if ms.lastBody["from_email"] != "from@test.com" {
		t.Errorf("from_email = %v", ms.lastBody["from_email"])
	}
	if ms.lastBody["api_key"] != "test-api-key" {
		t.Errorf("api_key = %v", ms.lastBody["api_key"])
	}
}

func TestSend_NilMessage(t *testing.T) {
	c := New("key")
	_, err := c.Send(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil message")
	}
}

func TestSend_EmptyTo(t *testing.T) {
	c := New("key")
	_, err := c.Send(context.Background(), &Message{Subject: "test"})
	if err == nil {
		t.Error("expected error for empty To")
	}
}

func TestSend_APIError(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, errResponse(4, "Invalid API key"))
	c := newTestClient(t, ms)

	_, err := c.Send(context.Background(), &Message{To: "user@test.com"})
	if err == nil {
		t.Fatal("expected API error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error is %T, want *APIError", err)
	}
	if apiErr.Code != 4 {
		t.Errorf("code = %d, want 4", apiErr.Code)
	}
	if apiErr.Message != "Invalid API key" {
		t.Errorf("message = %q", apiErr.Message)
	}
}

func TestSend_WithReplace(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{"transaction_id":"def456"}`))
	c := newTestClient(t, ms)

	_, err := c.Send(context.Background(), &Message{
		To:   "user@test.com",
		HTML: "<p>Hello %NAME%</p>",
		Replace: map[string]string{
			"%NAME%": "World",
		},
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	replace, ok := ms.lastBody["replace"].(map[string]any)
	if !ok {
		t.Fatalf("replace is %T", ms.lastBody["replace"])
	}
	if replace["%NAME%"] != "World" {
		t.Errorf("replace[%%NAME%%] = %v", replace["%NAME%"])
	}
}

func TestCheck_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{
		"date":"2024-01-15 10:00:00",
		"datesent":"2024-01-15 10:00:01",
		"to":"user@test.com",
		"status":5,
		"statusname":"Sent",
		"statuschangedate":"2024-01-15 10:00:02"
	}`))
	c := newTestClient(t, ms)

	resp, err := c.Check(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if resp.Status != 5 {
		t.Errorf("Status = %d, want 5", resp.Status)
	}
	if resp.StatusName != "Sent" {
		t.Errorf("StatusName = %q", resp.StatusName)
	}
	if resp.To != "user@test.com" {
		t.Errorf("To = %q", resp.To)
	}
	if resp.DateSent != "2024-01-15 10:00:01" {
		t.Errorf("DateSent = %q", resp.DateSent)
	}

	// Verify body
	if ms.lastBody["transaction_id"] != "abc123" {
		t.Errorf("transaction_id = %v", ms.lastBody["transaction_id"])
	}
}

func TestCheck_EmptyID(t *testing.T) {
	c := New("key")
	_, err := c.Check(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty transaction ID")
	}
}

func TestGetLog_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`[{"date":"2024-01-15","email":"user@test.com"}]`))
	c := newTestClient(t, ms)

	entries, err := c.GetLog(context.Background(), &GetLogParams{
		EventType: "delivered",
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("GetLog: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}

	if ms.lastBody["event_type"] != "delivered" {
		t.Errorf("event_type = %v", ms.lastBody["event_type"])
	}
}

func TestGetLog_NilParams(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`[]`))
	c := newTestClient(t, ms)

	entries, err := c.GetLog(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetLog: %v", err)
	}
	if entries == nil {
		t.Error("entries should not be nil")
	}
}

func TestGetLog_SingleObject(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{"date":"2024-01-15","email":"user@test.com"}`))
	c := newTestClient(t, ms)

	entries, err := c.GetLog(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetLog: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1 (single object wrapped)", len(entries))
	}
}

func TestGetStat_Success(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{"sent":100,"delivered":95}`))
	c := newTestClient(t, ms)

	data, err := c.GetStat(context.Background(), &GetStatParams{
		Period:    "custom",
		StartDate: "2024-01-01",
		EndDate:   "2024-01-31",
	})
	if err != nil {
		t.Fatalf("GetStat: %v", err)
	}
	if data == nil {
		t.Error("data should not be nil")
	}

	if ms.lastBody["period"] != "custom" {
		t.Errorf("period = %v", ms.lastBody["period"])
	}
	if ms.lastBody["start_date"] != "2024-01-01" {
		t.Errorf("start_date = %v", ms.lastBody["start_date"])
	}
}

func TestGetStat_NilParams(t *testing.T) {
	ms := newMockServer(t)
	ms.setResponse(200, okResponse(`{}`))
	c := newTestClient(t, ms)

	_, err := c.GetStat(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetStat: %v", err)
	}
}
