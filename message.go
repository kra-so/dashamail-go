package dashamail

import (
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

// Message describes an email to be sent via the transactional API.
type Message struct {
	// To is the recipient email address (required).
	To string `json:"to"`

	// Subject is the email subject line.
	Subject string `json:"subject,omitempty"`

	// HTML is the HTML body of the email.
	HTML string `json:"message,omitempty"`

	// PlainText is the plain-text fallback body.
	PlainText string `json:"plain_text,omitempty"`

	// FromEmail overrides the client-level sender address.
	FromEmail string `json:"from_email,omitempty"`

	// FromName overrides the client-level sender name.
	FromName string `json:"from_name,omitempty"`

	// CC is the carbon copy recipient(s).
	CC string `json:"cc,omitempty"`

	// BCC is the blind carbon copy recipient(s).
	BCC string `json:"bcc,omitempty"`

	// MessageID is a custom Message-ID header value.
	MessageID string `json:"message_id,omitempty"`

	// DeliveryTime schedules the email for a specific time (format: "YYYY-MM-DD HH:MM:SS").
	DeliveryTime string `json:"delivery_time,omitempty"`

	// Replace is a map of template tags to replacement values.
	// For example: {"%TAG1%": "value1", "%TAG2%": "value2"}
	Replace map[string]string `json:"replace,omitempty"`

	// Domain overrides the sending domain.
	Domain string `json:"domain,omitempty"`

	// Headers is a map of custom email headers.
	Headers map[string]string `json:"headers,omitempty"`

	// TemplateData is arbitrary data passed to the template engine.
	TemplateData map[string]any `json:"template_data,omitempty"`

	// NoTrackOpens overrides the client-level open tracking setting.
	// nil means use the client default.
	NoTrackOpens *bool `json:"no_track_opens,omitempty"`

	// NoTrackClicks overrides the client-level click tracking setting.
	// nil means use the client default.
	NoTrackClicks *bool `json:"no_track_clicks,omitempty"`

	// IgnoreDeliveryPolicy overrides the client-level delivery policy setting.
	// nil means use the client default.
	IgnoreDeliveryPolicy *bool `json:"ignore_delivery_policy,omitempty"`

	// Attachments is a list of file attachments.
	Attachments []Attachment `json:"attachments,omitempty"`

	// Inline is a list of inline images (referenced via cid: in HTML).
	Inline []InlineAttachment `json:"inline,omitempty"`
}

// Attachment represents a file attached to an email.
type Attachment struct {
	// Name is the filename as it will appear to the recipient.
	Name string `json:"name"`

	// FileBody is the Base64-encoded file content.
	FileBody string `json:"filebody"`
}

// InlineAttachment represents an inline image in an email.
type InlineAttachment struct {
	// MIMEType is the MIME type of the file (e.g. "image/png").
	MIMEType string `json:"mime_type"`

	// Filename is the original filename.
	Filename string `json:"filename"`

	// Body is the Base64-encoded file content.
	Body string `json:"body"`

	// CID is the Content-ID used to reference the image in HTML (e.g. <img src="cid:123">).
	CID string `json:"cid"`
}

// AttachFile reads a file from disk and appends it to msg.Attachments.
func (m *Message) AttachFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("dashamail: read attachment %q: %w", path, err)
	}
	m.Attachments = append(m.Attachments, Attachment{
		Name:     filepath.Base(path),
		FileBody: base64.StdEncoding.EncodeToString(data),
	})
	return nil
}

// AttachInlineFile reads a file from disk and appends it to msg.Inline with the given CID.
func (m *Message) AttachInlineFile(path, cid string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("dashamail: read inline attachment %q: %w", path, err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	m.Inline = append(m.Inline, InlineAttachment{
		MIMEType: mimeType,
		Filename: filepath.Base(path),
		Body:     base64.StdEncoding.EncodeToString(data),
		CID:      cid,
	})
	return nil
}

// toPayload converts the message to a map suitable for JSON serialization,
// merging client-level defaults where message-level overrides are not set.
func (m *Message) toPayload(c *Client) map[string]any {
	p := make(map[string]any)

	p["to"] = m.To

	if m.Subject != "" {
		p["subject"] = m.Subject
	}
	if m.HTML != "" {
		p["message"] = m.HTML
	}
	if m.PlainText != "" {
		p["plain_text"] = m.PlainText
	}

	// From: message-level overrides client-level.
	fromEmail := c.fromEmail
	if m.FromEmail != "" {
		fromEmail = m.FromEmail
	}
	if fromEmail != "" {
		p["from_email"] = fromEmail
	}

	fromName := c.fromName
	if m.FromName != "" {
		fromName = m.FromName
	}
	if fromName != "" {
		p["from_name"] = fromName
	}

	if m.CC != "" {
		p["cc"] = m.CC
	}
	if m.BCC != "" {
		p["bcc"] = m.BCC
	}
	if m.MessageID != "" {
		p["message_id"] = m.MessageID
	}
	if m.DeliveryTime != "" {
		p["delivery_time"] = m.DeliveryTime
	}
	if m.Domain != "" {
		p["domain"] = m.Domain
	}

	if len(m.Replace) > 0 {
		p["replace"] = m.Replace
	}
	if len(m.Headers) > 0 {
		p["headers"] = m.Headers
	}
	if len(m.TemplateData) > 0 {
		p["template_data"] = m.TemplateData
	}

	// Tracking flags: message-level overrides client-level.
	if m.NoTrackOpens != nil {
		p["no_track_opens"] = *m.NoTrackOpens
	} else {
		p["no_track_opens"] = c.noTrackOpens
	}
	if m.NoTrackClicks != nil {
		p["no_track_clicks"] = *m.NoTrackClicks
	} else {
		p["no_track_clicks"] = c.noTrackClicks
	}
	if m.IgnoreDeliveryPolicy != nil {
		p["ignore_delivery_policy"] = *m.IgnoreDeliveryPolicy
	} else {
		p["ignore_delivery_policy"] = c.ignoreDeliveryPolicy
	}

	if len(m.Attachments) > 0 {
		p["attachments"] = m.Attachments
	}
	if len(m.Inline) > 0 {
		p["inline"] = m.Inline
	}

	return p
}
