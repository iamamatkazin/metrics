package agent

import (
	"context"
	"log/slog"

	"github.com/iamamatkazin/metrics.git/cmd/crypto"
)

// Worker обрабатывает очередь задач отправки метрик на сервер.
// Читает задачи из канала jobs и отправляет их через HTTP клиент.
func (a *Agent) Worker() {
	for job := range a.jobs {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
			defer cancel()

			var (
				data []byte
				err  error
			)

			if job.metric != nil {
				if data, err = crypto.Encrypt(a.cfg.CryptoKey, job.metric); err != nil {
					slog.Error("Ошибка шифрования метрик:", slog.Any("error", err))
					return
				}
			}

			if err = a.client.Post(ctx, job.url, job.contentType, a.cfg.Key, data); err != nil {
				slog.Error("Ошибка отправки метрик на сервер:", slog.Any("error", err))
			}
		}()
	}
}
