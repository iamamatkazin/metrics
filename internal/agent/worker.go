package agent

import (
	"context"
	"log/slog"
)

func (a *Agent) Worker() {
	for job := range a.jobs {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
			defer cancel()

			if err := a.client.Post(ctx, job.url, job.contentType, job.metric); err != nil {
				slog.Error("Ошибка отправки метрик на сервер:", slog.Any("error", err))
			}
		}()
	}
}
