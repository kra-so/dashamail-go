package dashamail

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestMessage_ToPayload_Defaults(t *testing.T) {
	c := New("key",
		WithFromEmail("default@test.com"),
		WithFromName("Default Sender"),
	)

	msg := &Message{
		To:      "user@test.com",
		Subject: "Hello",
		HTML:    "<p>Hi</p>",
	}

	p := msg.toPayload(c)

	if p["to"] != "user@test.com" {
		t.Errorf("to = %v", p["to"])
	}
	if p["subject"] != "Hello" {
		t.Errorf("subject = %v", p["subject"])
	}
	if p["message"] != "<p>Hi</p>" {
		t.Errorf("message = %v", p["message"])
	}
	if p["from_email"] != "default@test.com" {
		t.Errorf("from_email = %v, want default@test.com", p["from_email"])
	}
	if p["from_name"] != "Default Sender" {
		t.Errorf("from_name = %v", p["from_name"])
	}
	if p["no_track_opens"] != true {
		t.Errorf("no_track_opens = %v, want true", p["no_track_opens"])
	}
	if p["no_track_clicks"] != true {
		t.Errorf("no_track_clicks = %v, want true", p["no_track_clicks"])
	}
	if p["ignore_delivery_policy"] != false {
		t.Errorf("ignore_delivery_policy = %v, want false", p["ignore_delivery_policy"])
	}
}

func TestMessage_ToPayload_Overrides(t *testing.T) {
	c := New("key",
		WithFromEmail("default@test.com"),
		WithFromName("Default"),
		WithNoTrackOpens(true),
	)

	msg := &Message{
		To:                   "user@test.com",
		FromEmail:            "override@test.com",
		FromName:             "Override",
		NoTrackOpens:         Bool(false),
		NoTrackClicks:        Bool(false),
		IgnoreDeliveryPolicy: Bool(true),
	}

	p := msg.toPayload(c)

	if p["from_email"] != "override@test.com" {
		t.Errorf("from_email = %v, want override", p["from_email"])
	}
	if p["from_name"] != "Override" {
		t.Errorf("from_name = %v", p["from_name"])
	}
	if p["no_track_opens"] != false {
		t.Errorf("no_track_opens = %v, want false", p["no_track_opens"])
	}
	if p["no_track_clicks"] != false {
		t.Errorf("no_track_clicks = %v, want false", p["no_track_clicks"])
	}
	if p["ignore_delivery_policy"] != true {
		t.Errorf("ignore_delivery_policy = %v, want true", p["ignore_delivery_policy"])
	}
}

func TestMessage_ToPayload_Replace(t *testing.T) {
	c := New("key")
	msg := &Message{
		To: "user@test.com",
		Replace: map[string]string{
			"%NAME%": "John",
		},
	}

	p := msg.toPayload(c)
	r, ok := p["replace"].(map[string]string)
	if !ok {
		t.Fatalf("replace is %T, want map[string]string", p["replace"])
	}
	if r["%NAME%"] != "John" {
		t.Errorf("replace[%%NAME%%] = %q", r["%NAME%"])
	}
}

func TestMessage_ToPayload_OptionalFields(t *testing.T) {
	c := New("key")
	msg := &Message{
		To:           "user@test.com",
		CC:           "cc@test.com",
		BCC:          "bcc@test.com",
		MessageID:    "custom-id",
		DeliveryTime: "2024-01-01 12:00:00",
		Domain:       "custom.com",
		PlainText:    "plain text body",
	}

	p := msg.toPayload(c)

	checks := map[string]string{
		"cc":            "cc@test.com",
		"bcc":           "bcc@test.com",
		"message_id":    "custom-id",
		"delivery_time": "2024-01-01 12:00:00",
		"domain":        "custom.com",
		"plain_text":    "plain text body",
	}
	for key, want := range checks {
		if got, ok := p[key].(string); !ok || got != want {
			t.Errorf("%s = %v, want %q", key, p[key], want)
		}
	}
}

func TestMessage_AttachFile(t *testing.T) {
	// Create a temp file
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := []byte("hello attachment")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	msg := &Message{To: "user@test.com"}
	if err := msg.AttachFile(path); err != nil {
		t.Fatal(err)
	}

	if len(msg.Attachments) != 1 {
		t.Fatalf("got %d attachments, want 1", len(msg.Attachments))
	}

	att := msg.Attachments[0]
	if att.Name != "test.txt" {
		t.Errorf("name = %q, want %q", att.Name, "test.txt")
	}

	decoded, err := base64.StdEncoding.DecodeString(att.FileBody)
	if err != nil {
		t.Fatalf("decode base64: %v", err)
	}
	if string(decoded) != "hello attachment" {
		t.Errorf("decoded = %q", string(decoded))
	}
}

func TestMessage_AttachFile_NotFound(t *testing.T) {
	msg := &Message{To: "user@test.com"}
	err := msg.AttachFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestMessage_AttachInlineFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logo.png")
	content := []byte("fake png data")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	msg := &Message{To: "user@test.com"}
	if err := msg.AttachInlineFile(path, "logo123"); err != nil {
		t.Fatal(err)
	}

	if len(msg.Inline) != 1 {
		t.Fatalf("got %d inline, want 1", len(msg.Inline))
	}

	inl := msg.Inline[0]
	if inl.Filename != "logo.png" {
		t.Errorf("filename = %q", inl.Filename)
	}
	if inl.CID != "logo123" {
		t.Errorf("cid = %q", inl.CID)
	}
	if inl.MIMEType != "image/png" {
		t.Errorf("mime_type = %q, want image/png", inl.MIMEType)
	}

	decoded, err := base64.StdEncoding.DecodeString(inl.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(decoded) != "fake png data" {
		t.Errorf("decoded = %q", string(decoded))
	}
}

func TestMessage_MultipleAttachments(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(name), 0644); err != nil {
			t.Fatal(err)
		}
	}

	msg := &Message{To: "user@test.com"}
	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := msg.AttachFile(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}

	if len(msg.Attachments) != 3 {
		t.Errorf("got %d attachments, want 3", len(msg.Attachments))
	}
}

func TestBool(t *testing.T) {
	tr := Bool(true)
	fl := Bool(false)

	if tr == nil || *tr != true {
		t.Error("Bool(true) failed")
	}
	if fl == nil || *fl != false {
		t.Error("Bool(false) failed")
	}
}
