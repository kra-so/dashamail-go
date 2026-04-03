package dashamail

import "encoding/json"

// RawResponse holds the raw API response data.
type RawResponse struct {
	// HTTPCode is the HTTP status code.
	HTTPCode int

	// Body is the raw response body bytes.
	Body []byte

	// Msg contains the API-level status message.
	Msg ResponseMsg

	// Data is the raw JSON of the "data" field in the response.
	Data json.RawMessage
}

// ResponseMsg represents the "msg" object in every DashaMail API response.
type ResponseMsg struct {
	ErrCode int    `json:"err_code"`
	Text    string `json:"text"`
	Type    string `json:"type"`
}

// SendResponse is returned by Client.Send.
type SendResponse struct {
	TransactionID string `json:"transaction_id"`
}

// CheckResponse is returned by Client.Check.
type CheckResponse struct {
	Date             string `json:"date"`
	DateSent         string `json:"datesent"`
	To               string `json:"to"`
	Status           int    `json:"status"`
	StatusName       string `json:"statusname"`
	StatusChangeDate string `json:"statuschangedate"`
}

// LogEntry represents a single entry from transactional.get_log.
type LogEntry struct {
	Date          string `json:"date"`
	Email         string `json:"email"`
	Subject       string `json:"subject"`
	Event         string `json:"event"`
	TransactionID string `json:"transaction_id"`
}

// GetLogResponse is returned by Client.GetLog.
type GetLogResponse struct {
	Entries []LogEntry
}

// StatEntry represents a statistics record from transactional.get_stat.
type StatEntry struct {
	Date      string `json:"date"`
	Sent      int    `json:"sent"`
	Delivered int    `json:"delivered"`
	Opened    int    `json:"opened"`
	Clicked   int    `json:"clicked"`
	Bounced   int    `json:"bounced"`
	Spam      int    `json:"spam"`
	Unsub     int    `json:"unsub"`
}

// WebhookURLs holds the current transactional webhook URLs.
type WebhookURLs struct {
	Open      string `json:"open,omitempty"`
	Click     string `json:"click,omitempty"`
	Hard      string `json:"hard,omitempty"`
	Soft      string `json:"soft,omitempty"`
	Spam      string `json:"spam,omitempty"`
	Unsub     string `json:"unsub,omitempty"`
	Subscribe string `json:"subscribe,omitempty"`
	Confirm   string `json:"confirm,omitempty"`
}
