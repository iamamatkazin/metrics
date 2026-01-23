package filestorage

import (
	"encoding/json"
	"io"
	"os"

	"github.com/iamamatkazin/metrics.git/internal/model"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Storage - структура, реализующая интерфейс работы с метриками для файлового хранилища.
type Storage struct {
	cfg  *sconfig.Config
	file *os.File
}

// New - конструктор для Storage.
func New(cfg *sconfig.Config) (*Storage, error) {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		cfg:  cfg,
		file: file,
	}

	return s, nil
}

// Close - завершить работу с файловым хранилищем.
func (s *Storage) Close() {
	if s.file != nil {
		s.file.Close()
	}
}

// SaveToFile - сохранить метрики в файл.
func (s *Storage) SaveToFile(metrics map[string]*model.Metric) error {
	if err := s.file.Truncate(0); err != nil {
		return err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	if err := json.NewEncoder(s.file).Encode(metrics); err != nil {
		return err
	}

	return nil
}

// LoadDump - загрузить сохраненные метрики в память.
func (s *Storage) LoadDump(metrics map[string]*model.Metric) error {
	if err := json.NewDecoder(s.file).Decode(&metrics); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}
