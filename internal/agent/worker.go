package agent

import (
	"context"
	"log/slog"
)

// worker это наш рабочий, который принимает два канала:
// jobs - канал задач, это входные данные для обработки
func (a *Agent) worker(ctx context.Context, jobs <-chan request) {
	for {
		select {
		case <-ctx.Done():
			return

		case job := <-jobs:
			if err := a.client.Post(ctx, job.url, job.contentType, job.metric); err != nil {
				slog.Error("Ошибка отправки метрик на сервер:", slog.Any("error", err))
			}
		}
	}
}
