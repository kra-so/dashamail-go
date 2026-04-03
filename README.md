# dashamail-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kra-so/dashamail-go.svg)](https://pkg.go.dev/github.com/kra-so/dashamail-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Go client library for the [DashaMail transactional email API](https://dashamail.ru/transactional/).

**[Документация на русском](README.ru.md)**

## Features

- Sending transactional emails with HTML, plain text, attachments, and inline images
- Template tag substitution (`replace`)
- Delivery status checking
- Event logs and statistics
- Transactional webhook management
- Functional options for client configuration
- Per-message overrides for sender, tracking, and delivery policy
- Typed API errors with `errors.As` support

## Installation

```bash
go get github.com/kra-so/dashamail-go
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	dashamail "github.com/kra-so/dashamail-go"
)

func main() {
	client := dashamail.New("your-api-key",
		dashamail.WithFromEmail("noreply@example.com"),
		dashamail.WithFromName("My App"),
	)

	resp, err := client.Send(context.Background(), &dashamail.Message{
		To:      "user@example.com",
		Subject: "Welcome!",
		HTML:    "<h1>Hello!</h1><p>Welcome to our service.</p>",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction ID:", resp.TransactionID)
}
```

## Client Configuration

```go
client := dashamail.New("api-key",
	dashamail.WithEndpoint("https://api.dashamail.com"),  // default
	dashamail.WithFromEmail("noreply@example.com"),       // sender email
	dashamail.WithFromName("My App"),                     // sender display name
	dashamail.WithNoTrackOpens(true),                     // disable open tracking (default: true)
	dashamail.WithNoTrackClicks(true),                    // disable click tracking (default: true)
	dashamail.WithIgnoreDeliveryPolicy(false),            // ignore delivery policy (default: false)
	dashamail.WithHTTPClient(customHTTPClient),           // custom *http.Client
	dashamail.WithDebug(true),                            // enable debug mode
)
```

## Sending Emails

### Simple Email

```go
resp, err := client.Send(ctx, &dashamail.Message{
	To:      "user@example.com",
	Subject: "Email subject",
	HTML:    "<p>HTML content</p>",
})
```

### With Template Substitution

```go
resp, err := client.Send(ctx, &dashamail.Message{
	To:      "user@example.com",
	Subject: "Your order",
	HTML:    "<p>Hello, %NAME%! Your order #%ORDER% has been received.</p>",
	Replace: map[string]string{
		"%NAME%":  "John",
		"%ORDER%": "12345",
	},
})
```

### With Attachments

```go
msg := &dashamail.Message{
	To:      "user@example.com",
	Subject: "Report",
	HTML:    `<p>Report attached.</p><img src="cid:logo">`,
}

// File attachment
if err := msg.AttachFile("./report.pdf"); err != nil {
	log.Fatal(err)
}

// Inline image (referenced via cid: in HTML)
if err := msg.AttachInlineFile("./logo.png", "logo"); err != nil {
	log.Fatal(err)
}

resp, err := client.Send(ctx, msg)
```

### Per-Message Overrides

Any client-level default can be overridden on an individual message:

```go
resp, err := client.Send(ctx, &dashamail.Message{
	To:           "user@example.com",
	Subject:      "Urgent",
	HTML:         "<p>Time-sensitive email</p>",
	FromEmail:    "urgent@example.com",
	FromName:     "Urgent Bot",
	NoTrackOpens: dashamail.Bool(false), // enable open tracking for this message
})
```

## Checking Delivery Status

```go
status, err := client.Check(ctx, "5a802b10ba82eccfd164f3c8be0fb678")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Status: %s (%d)\n", status.StatusName, status.Status)
fmt.Printf("Sent at: %s\n", status.DateSent)
```

## Event Logs

```go
entries, err := client.GetLog(ctx, &dashamail.GetLogParams{
	EventType: "delivered",
	From:      "2024-01-01 00:00:00",
	To:        "2024-01-31 23:59:59",
	Limit:     100,
})
```

## Statistics

```go
data, err := client.GetStat(ctx, &dashamail.GetStatParams{
	Period:    "custom",
	StartDate: "2024-01-01",
	EndDate:   "2024-01-31",
})
```

## Webhooks

```go
// Set webhook URLs
err := client.SetTransactionalWebhooks(ctx, &dashamail.WebhookURLs{
	Open:  "https://example.com/webhooks/open",
	Click: "https://example.com/webhooks/click",
	Hard:  "https://example.com/webhooks/bounce",
})

// Get current webhooks
data, err := client.GetTransactionalWebhooks(ctx, "")

// Delete a webhook
err := client.DeleteTransactionalWebhooks(ctx, "open")
```

## Error Handling

API errors are returned as `*APIError` and can be inspected using `errors.As`:

```go
resp, err := client.Send(ctx, msg)
if err != nil {
	var apiErr *dashamail.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}
```

## Message Fields Reference

| Field | Type | Description |
|---|---|---|
| `To` | `string` | Recipient email (required) |
| `Subject` | `string` | Email subject |
| `HTML` | `string` | HTML body |
| `PlainText` | `string` | Plain-text fallback body |
| `FromEmail` | `string` | Sender email override |
| `FromName` | `string` | Sender name override |
| `CC` | `string` | Carbon copy recipient(s) |
| `BCC` | `string` | Blind carbon copy recipient(s) |
| `MessageID` | `string` | Custom Message-ID header |
| `DeliveryTime` | `string` | Scheduled delivery time (`YYYY-MM-DD HH:MM:SS`) |
| `Replace` | `map[string]string` | Template tag substitutions |
| `Domain` | `string` | Sending domain override |
| `Headers` | `map[string]string` | Custom email headers |
| `TemplateData` | `map[string]any` | Data for template engine |
| `NoTrackOpens` | `*bool` | Open tracking override |
| `NoTrackClicks` | `*bool` | Click tracking override |
| `IgnoreDeliveryPolicy` | `*bool` | Delivery policy override |
| `Attachments` | `[]Attachment` | File attachments |
| `Inline` | `[]InlineAttachment` | Inline images |

## Testing

```bash
go test ./...
```

## License

MIT
