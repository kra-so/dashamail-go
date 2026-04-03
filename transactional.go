package dashamail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Send sends a transactional email.
//
// On success it returns the transaction ID which can be used with Check.
func (c *Client) Send(ctx context.Context, msg *Message) (*SendResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("dashamail: message must not be nil")
	}
	if msg.To == "" {
		return nil, fmt.Errorf("dashamail: message.To is required")
	}

	payload := msg.toPayload(c)

	raw, err := c.do(ctx, "transactional.send", http.MethodPost, payload)
	if err != nil {
		return nil, err
	}

	var resp SendResponse
	if err := json.Unmarshal(raw.Data, &resp); err != nil {
		return nil, fmt.Errorf("dashamail: decode send response: %w", err)
	}
	return &resp, nil
}

// Check returns the delivery status of a previously sent transactional email.
func (c *Client) Check(ctx context.Context, transactionID string) (*CheckResponse, error) {
	if transactionID == "" {
		return nil, fmt.Errorf("dashamail: transactionID is required")
	}

	payload := map[string]any{
		"transaction_id": transactionID,
	}

	raw, err := c.do(ctx, "transactional.check", http.MethodPost, payload)
	if err != nil {
		return nil, err
	}

	var resp CheckResponse
	if err := json.Unmarshal(raw.Data, &resp); err != nil {
		return nil, fmt.Errorf("dashamail: decode check response: %w", err)
	}
	return &resp, nil
}

// GetLogParams configures the transactional.get_log request.
type GetLogParams struct {
	// EventType filters by event type (e.g. "delivered", "opened", "clicked", "bounced", "spam").
	EventType string `json:"event_type,omitempty"`

	// Emails filters by recipient emails.
	Emails []string `json:"emails,omitempty"`

	// Sort specifies the sort order ("asc" or "desc").
	Sort string `json:"sort,omitempty"`

	// CampaignID filters by campaign.
	CampaignID string `json:"campaign_id,omitempty"`

	// Start is the offset for pagination.
	Start int `json:"start,omitempty"`

	// Limit is the maximum number of entries to return.
	Limit int `json:"limit,omitempty"`

	// From filters events from this time (format: "YYYY-MM-DD HH:MM:SS").
	From string `json:"from,omitempty"`

	// To filters events until this time (format: "YYYY-MM-DD HH:MM:SS").
	To string `json:"to,omitempty"`
}

// GetLog returns transactional email event logs.
func (c *Client) GetLog(ctx context.Context, params *GetLogParams) ([]json.RawMessage, error) {
	var payload map[string]any
	if params != nil {
		raw, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("dashamail: marshal get_log params: %w", err)
		}
		payload = make(map[string]any)
		if err := json.Unmarshal(raw, &payload); err != nil {
			return nil, fmt.Errorf("dashamail: unmarshal get_log params: %w", err)
		}
	} else {
		payload = make(map[string]any)
	}

	resp, err := c.do(ctx, "transactional.get_log", http.MethodPost, payload)
	if err != nil {
		return nil, err
	}

	// The data field can be an array or object; return raw for flexibility.
	var entries []json.RawMessage
	if len(resp.Data) > 0 && resp.Data[0] == '[' {
		if err := json.Unmarshal(resp.Data, &entries); err != nil {
			return nil, fmt.Errorf("dashamail: decode get_log response: %w", err)
		}
	} else {
		// Single object — wrap it.
		entries = append(entries, resp.Data)
	}
	return entries, nil
}

// GetStatParams configures the transactional.get_stat request.
type GetStatParams struct {
	// Period is the statistics period: "today", "yesterday", "last_7_days",
	// "last_30_days", "last_90_days", "custom".
	Period string `json:"period,omitempty"`

	// StartDate is used with Period="custom" (format: "YYYY-MM-DD").
	StartDate string `json:"start_date,omitempty"`

	// EndDate is used with Period="custom" (format: "YYYY-MM-DD").
	EndDate string `json:"end_date,omitempty"`
}

// GetStat returns transactional email statistics for a time period.
func (c *Client) GetStat(ctx context.Context, params *GetStatParams) (json.RawMessage, error) {
	var payload map[string]any
	if params != nil {
		raw, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("dashamail: marshal get_stat params: %w", err)
		}
		payload = make(map[string]any)
		if err := json.Unmarshal(raw, &payload); err != nil {
			return nil, fmt.Errorf("dashamail: unmarshal get_stat params: %w", err)
		}
	} else {
		payload = make(map[string]any)
	}

	resp, err := c.do(ctx, "transactional.get_stat", http.MethodPost, payload)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
