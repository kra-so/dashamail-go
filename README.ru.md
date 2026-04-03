# dashamail-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kra-so/dashamail-go.svg)](https://pkg.go.dev/github.com/kra-so/dashamail-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Go-клиент для [транзакционного email API DashaMail](https://dashamail.ru/transactional/).

**[English documentation](README.md)**

## Возможности

- Отправка транзакционных писем с HTML, plain text, вложениями и inline-изображениями
- Подстановка шаблонных тегов (`replace`)
- Проверка статуса доставки
- Получение логов событий и статистики
- Управление транзакционными вебхуками
- Конфигурация клиента через функциональные опции
- Переопределение параметров на уровне каждого письма
- Типизированные ошибки API с поддержкой `errors.As`

## Установка

```bash
go get github.com/kra-so/dashamail-go
```

Требуется Go 1.21 или выше.

## Быстрый старт

```go
package main

import (
	"context"
	"fmt"
	"log"

	dashamail "github.com/kra-so/dashamail-go"
)

func main() {
	client := dashamail.New("ваш-api-ключ",
		dashamail.WithFromEmail("noreply@example.com"),
		dashamail.WithFromName("Моё приложение"),
	)

	resp, err := client.Send(context.Background(), &dashamail.Message{
		To:      "user@example.com",
		Subject: "Добро пожаловать!",
		HTML:    "<h1>Привет!</h1><p>Добро пожаловать в наш сервис.</p>",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction ID:", resp.TransactionID)
}
```

## Конфигурация клиента

```go
client := dashamail.New("api-key",
	dashamail.WithEndpoint("https://api.dashamail.com"),  // по умолчанию
	dashamail.WithFromEmail("noreply@example.com"),       // email отправителя
	dashamail.WithFromName("My App"),                     // имя отправителя
	dashamail.WithNoTrackOpens(true),                     // не отслеживать открытия (по умолчанию: true)
	dashamail.WithNoTrackClicks(true),                    // не отслеживать клики (по умолчанию: true)
	dashamail.WithIgnoreDeliveryPolicy(false),            // игнорировать политику доставки (по умолчанию: false)
	dashamail.WithHTTPClient(customHTTPClient),           // свой *http.Client
	dashamail.WithDebug(true),                            // режим отладки
)
```

## Отправка писем

### Простое письмо

```go
resp, err := client.Send(ctx, &dashamail.Message{
	To:      "user@example.com",
	Subject: "Тема письма",
	HTML:    "<p>HTML-содержимое</p>",
})
```

### С шаблонными подстановками

```go
resp, err := client.Send(ctx, &dashamail.Message{
	To:      "user@example.com",
	Subject: "Ваш заказ",
	HTML:    "<p>Здравствуйте, %NAME%! Ваш заказ #%ORDER% принят.</p>",
	Replace: map[string]string{
		"%NAME%":  "Иван",
		"%ORDER%": "12345",
	},
})
```

### С вложениями

```go
msg := &dashamail.Message{
	To:      "user@example.com",
	Subject: "Отчёт",
	HTML:    `<p>Отчёт во вложении.</p><img src="cid:logo">`,
}

// Файловое вложение
if err := msg.AttachFile("./report.pdf"); err != nil {
	log.Fatal(err)
}

// Inline-изображение (для ссылок cid: в HTML)
if err := msg.AttachInlineFile("./logo.png", "logo"); err != nil {
	log.Fatal(err)
}

resp, err := client.Send(ctx, msg)
```

### Переопределение параметров на уровне письма

Любой параметр, заданный на уровне клиента, можно переопределить для конкретного письма:

```go
resp, err := client.Send(ctx, &dashamail.Message{
	To:           "user@example.com",
	Subject:      "Важное",
	HTML:         "<p>Срочное письмо</p>",
	FromEmail:    "urgent@example.com",
	FromName:     "Urgent Bot",
	NoTrackOpens: dashamail.Bool(false), // включить отслеживание открытий для этого письма
})
```

## Проверка статуса доставки

```go
status, err := client.Check(ctx, "5a802b10ba82eccfd164f3c8be0fb678")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Статус: %s (%d)\n", status.StatusName, status.Status)
fmt.Printf("Отправлено: %s\n", status.DateSent)
```

## Логи событий

```go
entries, err := client.GetLog(ctx, &dashamail.GetLogParams{
	EventType: "delivered",
	From:      "2024-01-01 00:00:00",
	To:        "2024-01-31 23:59:59",
	Limit:     100,
})
```

## Статистика

```go
data, err := client.GetStat(ctx, &dashamail.GetStatParams{
	Period:    "custom",
	StartDate: "2024-01-01",
	EndDate:   "2024-01-31",
})
```

## Вебхуки

```go
// Установить URL-ы вебхуков
err := client.SetTransactionalWebhooks(ctx, &dashamail.WebhookURLs{
	Open:  "https://example.com/webhooks/open",
	Click: "https://example.com/webhooks/click",
	Hard:  "https://example.com/webhooks/bounce",
})

// Получить текущие вебхуки
data, err := client.GetTransactionalWebhooks(ctx, "")

// Удалить вебхук
err := client.DeleteTransactionalWebhooks(ctx, "open")
```

## Обработка ошибок

Ошибки API возвращаются как `*APIError` и могут быть проверены через `errors.As`:

```go
resp, err := client.Send(ctx, msg)
if err != nil {
	var apiErr *dashamail.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("Ошибка API %d: %s\n", apiErr.Code, apiErr.Message)
	} else {
		fmt.Printf("Ошибка: %v\n", err)
	}
}
```

## Справочник полей Message

| Поле | Тип | Описание |
|---|---|---|
| `To` | `string` | Email получателя (обязательно) |
| `Subject` | `string` | Тема письма |
| `HTML` | `string` | HTML-тело письма |
| `PlainText` | `string` | Текстовое тело (fallback) |
| `FromEmail` | `string` | Email отправителя (переопределение) |
| `FromName` | `string` | Имя отправителя (переопределение) |
| `CC` | `string` | Копия |
| `BCC` | `string` | Скрытая копия |
| `MessageID` | `string` | Пользовательский заголовок Message-ID |
| `DeliveryTime` | `string` | Время отложенной отправки (`YYYY-MM-DD HH:MM:SS`) |
| `Replace` | `map[string]string` | Подстановки шаблонных тегов |
| `Domain` | `string` | Домен отправки (переопределение) |
| `Headers` | `map[string]string` | Пользовательские заголовки |
| `TemplateData` | `map[string]any` | Данные для шаблонизатора |
| `NoTrackOpens` | `*bool` | Отслеживание открытий (переопределение) |
| `NoTrackClicks` | `*bool` | Отслеживание кликов (переопределение) |
| `IgnoreDeliveryPolicy` | `*bool` | Политика доставки (переопределение) |
| `Attachments` | `[]Attachment` | Файловые вложения |
| `Inline` | `[]InlineAttachment` | Inline-изображения |

## Запуск тестов

```bash
go test ./...
```

## Лицензия

MIT
