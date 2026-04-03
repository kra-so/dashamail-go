// Package dashamail provides a Go client for the DashaMail transactional email API.
//
// Usage:
//
//	client := dashamail.New("your-api-key",
//		dashamail.WithFromEmail("noreply@example.com"),
//		dashamail.WithFromName("My App"),
//	)
//
//	resp, err := client.Send(ctx, &dashamail.Message{
//		To:      "user@example.com",
//		Subject: "Hello",
//		HTML:    "<h1>Hi!</h1>",
//	})
package dashamail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultEndpoint is the default DashaMail API base URL.
	DefaultEndpoint = "https://api.dashamail.com"

	// Version is the library version.
	Version = "0.1.0"

	userAgent = "DashaMail(Go)/" + Version
)

// Client is the DashaMail API client.
type Client struct {
	apiKey   string
	endpoint string

	fromEmail string
	fromName  string

	noTrackOpens          bool
	noTrackClicks         bool
	ignoreDeliveryPolicy  bool

	httpClient *http.Client
	debug      bool
}

// Option configures the Client.
type Option func(*Client)

// New creates a new DashaMail client with the given API key and options.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:        apiKey,
		endpoint:      DefaultEndpoint,
		noTrackOpens:  true,
		noTrackClicks: true,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithEndpoint sets a custom API endpoint.
func WithEndpoint(endpoint string) Option {
	return func(c *Client) { c.endpoint = endpoint }
}

// WithFromEmail sets the default sender email address.
func WithFromEmail(email string) Option {
	return func(c *Client) { c.fromEmail = email }
}

// WithFromName sets the default sender display name.
func WithFromName(name string) Option {
	return func(c *Client) { c.fromName = name }
}

// WithNoTrackOpens disables or enables open tracking (default: disabled).
func WithNoTrackOpens(v bool) Option {
	return func(c *Client) { c.noTrackOpens = v }
}

// WithNoTrackClicks disables or enables click tracking (default: disabled).
func WithNoTrackClicks(v bool) Option {
	return func(c *Client) { c.noTrackClicks = v }
}

// WithIgnoreDeliveryPolicy sets whether to ignore the delivery policy.
func WithIgnoreDeliveryPolicy(v bool) Option {
	return func(c *Client) { c.ignoreDeliveryPolicy = v }
}

// WithHTTPClient sets a custom *http.Client for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithDebug enables debug logging to stderr.
func WithDebug(v bool) Option {
	return func(c *Client) { c.debug = v }
}

// do executes an API request. The method parameter is the DashaMail API method
// name (e.g. "transactional.send"), httpMethod is "GET" or "POST", and body
// is the JSON-serializable request payload (may be nil for GET requests).
func (c *Client) do(ctx context.Context, apiMethod, httpMethod string, body any) (*RawResponse, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("dashamail: invalid endpoint: %w", err)
	}

	q := u.Query()
	q.Set("method", apiMethod)
	q.Set("format", "JSON")
	u.RawQuery = q.Encode()

	var reqBody io.Reader
	if body != nil {
		payload, ok := body.(map[string]any)
		if !ok {
			// Marshal and re-unmarshal to get map
			raw, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("dashamail: marshal request: %w", err)
			}
			payload = make(map[string]any)
			if err := json.Unmarshal(raw, &payload); err != nil {
				return nil, fmt.Errorf("dashamail: unmarshal request: %w", err)
			}
		}
		payload["api_key"] = c.apiKey

		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("dashamail: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("dashamail: create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dashamail: http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("dashamail: read response: %w", err)
	}

	var raw RawResponse
	raw.HTTPCode = resp.StatusCode
	raw.Body = respBody

	var envelope struct {
		Response struct {
			Msg  ResponseMsg    `json:"msg"`
			Data json.RawMessage `json:"data"`
		} `json:"response"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("dashamail: decode response: %w (body: %s)", err, string(respBody))
	}

	raw.Msg = envelope.Response.Msg
	raw.Data = envelope.Response.Data

	if raw.Msg.ErrCode != 0 {
		return &raw, &APIError{
			Code:    raw.Msg.ErrCode,
			Message: raw.Msg.Text,
		}
	}

	return &raw, nil
}
