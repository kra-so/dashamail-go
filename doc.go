// Package dashamail provides a Go client for the DashaMail transactional email API
// (https://dashamail.ru/transactional/).
//
// # Quick Start
//
// Create a client and send a transactional email:
//
//	client := dashamail.New("your-api-key",
//		dashamail.WithFromEmail("noreply@example.com"),
//		dashamail.WithFromName("My App"),
//	)
//
//	resp, err := client.Send(ctx, &dashamail.Message{
//		To:      "user@example.com",
//		Subject: "Welcome!",
//		HTML:    "<h1>Hello!</h1>",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Transaction ID:", resp.TransactionID)
//
// # Configuration
//
// The client is configured via functional options passed to [New]:
//
//   - [WithEndpoint] — custom API base URL (default: https://api.dashamail.com)
//   - [WithFromEmail] — default sender email address
//   - [WithFromName] — default sender display name
//   - [WithNoTrackOpens] — disable open tracking (default: true)
//   - [WithNoTrackClicks] — disable click tracking (default: true)
//   - [WithIgnoreDeliveryPolicy] — ignore delivery policy (default: false)
//   - [WithHTTPClient] — custom *http.Client
//   - [WithDebug] — enable debug mode
//
// Per-message overrides are available via [Message] fields. Use [Bool] helper
// to set optional boolean fields.
//
// # Transactional API
//
// The following transactional API methods are supported:
//
//   - [Client.Send] — send a transactional email (transactional.send)
//   - [Client.Check] — check delivery status (transactional.check)
//   - [Client.GetLog] — retrieve event logs (transactional.get_log)
//   - [Client.GetStat] — retrieve statistics (transactional.get_stat)
//
// # Webhooks
//
//   - [Client.GetTransactionalWebhooks] — get configured webhook URLs
//   - [Client.SetTransactionalWebhooks] — set webhook URLs
//   - [Client.DeleteTransactionalWebhooks] — remove a webhook
//
// # Attachments
//
// Use [Message.AttachFile] and [Message.AttachInlineFile] to add attachments
// from disk:
//
//	msg := &dashamail.Message{To: "user@example.com", Subject: "Report"}
//	msg.AttachFile("report.pdf")
//	msg.AttachInlineFile("logo.png", "logo-cid")
//
// # Error Handling
//
// API errors are returned as [*APIError] and can be inspected with errors.As:
//
//	var apiErr *dashamail.APIError
//	if errors.As(err, &apiErr) {
//		fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
//	}
//
// ---
//
// Пакет dashamail — Go-клиент для транзакционного email API DashaMail
// (https://dashamail.ru/transactional/).
//
// # Быстрый старт
//
// Создайте клиент и отправьте транзакционное письмо:
//
//	client := dashamail.New("ваш-api-ключ",
//		dashamail.WithFromEmail("noreply@example.com"),
//		dashamail.WithFromName("Моё приложение"),
//	)
//
//	resp, err := client.Send(ctx, &dashamail.Message{
//		To:      "user@example.com",
//		Subject: "Добро пожаловать!",
//		HTML:    "<h1>Привет!</h1>",
//	})
//
// # Конфигурация
//
// Клиент настраивается через функциональные опции, передаваемые в [New].
// Каждое сообщение может переопределить настройки клиента через поля [Message].
// Используйте хелпер [Bool] для optional-булевых полей.
//
// # Транзакционное API
//
//   - [Client.Send] — отправка письма (transactional.send)
//   - [Client.Check] — проверка статуса доставки (transactional.check)
//   - [Client.GetLog] — получение логов событий (transactional.get_log)
//   - [Client.GetStat] — получение статистики (transactional.get_stat)
//
// # Вебхуки
//
//   - [Client.GetTransactionalWebhooks] — получить URL-ы вебхуков
//   - [Client.SetTransactionalWebhooks] — установить вебхуки
//   - [Client.DeleteTransactionalWebhooks] — удалить вебхук
//
// # Обработка ошибок
//
// Ошибки API возвращаются как [*APIError]. Проверяйте через errors.As.
package dashamail
