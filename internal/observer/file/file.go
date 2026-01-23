package file

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/iamamatkazin/metrics.git/internal/model"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Subscriber - структура файлового аудита.
type Subscriber struct {
	file *os.File
	id   string
}

// New - конструктор файлового аудита.
func New(cfg *sconfig.Config, id string) (*Subscriber, error) {
	file, err := os.OpenFile(cfg.FileAudit, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		id:   id,
		file: file,
	}, nil
}

// Send - отослать данные аудита в файл.
func (s *Subscriber) Send(mes model.Message) {
	if err := json.NewEncoder(s.file).Encode(mes); err != nil {
		slog.Error("Ошибка записи логов в лог файл:", slog.Any("error", err))
	}
}

// GetID - идентификатор файлового аудита.
func (s *Subscriber) GetID() string {
	return s.id
}
