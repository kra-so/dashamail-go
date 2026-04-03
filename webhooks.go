package dashamail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetTransactionalWebhooks returns the currently configured transactional webhook URLs.
//
// If eventName is not empty, only the specified event webhook is returned.
func (c *Client) GetTransactionalWebhooks(ctx context.Context, eventName string) (json.RawMessage, error) {
	payload := map[string]any{}
	if eventName != "" {
		payload["event_name"] = eventName
	}

	resp, err := c.do(ctx, "account.get_tr_webhooks", http.MethodPost, payload)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// SetTransactionalWebhooks sets the transactional webhook URLs.
//
// Pass only the event URLs you want to configure; empty strings are omitted.
func (c *Client) SetTransactionalWebhooks(ctx context.Context, urls *WebhookURLs) error {
	if urls == nil {
		return fmt.Errorf("dashamail: webhook URLs must not be nil")
	}

	payload := make(map[string]any)
	if urls.Open != "" {
		payload["open"] = urls.Open
	}
	if urls.Click != "" {
		payload["click"] = urls.Click
	}
	if urls.Hard != "" {
		payload["hard"] = urls.Hard
	}
	if urls.Soft != "" {
		payload["soft"] = urls.Soft
	}
	if urls.Spam != "" {
		payload["spam"] = urls.Spam
	}
	if urls.Unsub != "" {
		payload["unsub"] = urls.Unsub
	}
	if urls.Subscribe != "" {
		payload["subscribe"] = urls.Subscribe
	}
	if urls.Confirm != "" {
		payload["confirm"] = urls.Confirm
	}

	_, err := c.do(ctx, "account.add_tr_webhooks", http.MethodPost, payload)
	return err
}

// DeleteTransactionalWebhooks deletes a transactional webhook by event name.
//
// Event names: "open", "click", "hard", "soft", "spam", "unsub", "subscribe", "confirm".
func (c *Client) DeleteTransactionalWebhooks(ctx context.Context, eventName string) error {
	if eventName == "" {
		return fmt.Errorf("dashamail: eventName is required")
	}

	payload := map[string]any{
		"event_name": eventName,
	}

	_, err := c.do(ctx, "account.delete_tr_webhooks", http.MethodPost, payload)
	return err
}
