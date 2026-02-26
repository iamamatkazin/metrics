// Package url предоставляет HTTP-подписчик аудита для отправки изменений метрик.
// Реализует интерфейс Observer для отправки аудит-сообщений во внешние
// HTTP эндпоинты.
package url

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/iamamatkazin/metrics.git/internal/model"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
	pkghttp "github.com/iamamatkazin/metrics.git/pkg/http"
)

// Subscriber - структура внешнего аудита.
type Subscriber struct {
	cfg    *sconfig.Config
	client pkghttp.Clienter
	id     string
}

// New создает новый экземпляр внешнего аудита.
func New(cfg *sconfig.Config, id string) *Subscriber {
	return &Subscriber{
		id:     id,
		cfg:    cfg,
		client: pkghttp.New(cfg.Timeout),
	}
}

// Send отправляет данные аудита во внешний сервис.
func (s *Subscriber) Send(mes model.Message) {
	byteMes, err := json.Marshal(mes)
	if err != nil {
		slog.Error("Ошибка подготовки логов для отправки в сервис логирования:", slog.Any("error", err))
		return
	}

	if err := s.client.Post(context.Background(), s.cfg.URLAudit, "application/json", "", byteMes); err != nil {
		slog.Error("Ошибка отправки логов в сервис логирования:", slog.Any("error", err))
	}
}

// GetID возвращает идентификатор внешнего сервиса аудита.
func (s *Subscriber) GetID() string {
	return s.id
}
