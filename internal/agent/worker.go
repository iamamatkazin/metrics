package agent

import (
	"context"
	"log/slog"
)

// Worker обрабатывает очередь задач отправки метрик на сервер.
// Читает задачи из канала jobs и отправляет их через HTTP клиент.
func (a *Agent) Worker() {
	for job := range a.jobs {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
			defer cancel()

			if err := a.client.Post(ctx, job.url, job.contentType, a.cfg.Key, job.metric); err != nil {
				slog.Error("Ошибка отправки метрик на сервер:", slog.Any("error", err))
			}
		}()
	}
}
