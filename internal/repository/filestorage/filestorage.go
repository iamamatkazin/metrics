package filestorage

import (
	"encoding/json"
	"io"
	"os"

	"github.com/iamamatkazin/metrics.git/internal/model"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type Storage struct {
	cfg  *sconfig.Config
	file *os.File
}

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

func (s *Storage) Close() {
	if s.file != nil {
		s.file.Close()
	}
}

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

func (s *Storage) LoadDump(metrics map[string]*model.Metric) error {
	if err := json.NewDecoder(s.file).Decode(&metrics); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}
