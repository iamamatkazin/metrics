package url

import (
	"context"
	"log/slog"

	"github.com/iamamatkazin/metrics.git/internal/model"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
	pkghttp "github.com/iamamatkazin/metrics.git/pkg/http"
)

// реализация observer
type Subscriber struct {
	cfg    *sconfig.Config
	client pkghttp.Clienter
	id     string
}

func New(cfg *sconfig.Config, id string) *Subscriber {
	return &Subscriber{
		id:     id,
		cfg:    cfg,
		client: pkghttp.New(cfg.Timeout),
	}
}

func (s *Subscriber) Send(mes model.Message) {
	if err := s.client.Post(context.Background(), s.cfg.URLAudit, "application/json", "", mes); err != nil {
		slog.Error("Ошибка отправки логов в сервис логирования:", slog.Any("error", err))
	}
}

func (s *Subscriber) GetID() string {
	return s.id
}
