package model

import (
	"fmt"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Ограничиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// чтобы отличать значение "0" от незаданного значения
// и соответственно не кодировать в структуру.
// generate:reset
type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int     `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"-"`
}

// Validate выполняет валидацию структуры метрик.
func (m *Metric) Validate() error {
	if m.MType != Gauge && m.MType != Counter {
		return fmt.Errorf("неизвестный тип метрики: %s", m.MType)
	}

	return nil
}

// ValidateJSON выполняет валидацию структуры метрик в формате JSON.
func (m *Metric) ValidateJSON() error {
	if m.MType != Gauge && m.MType != Counter {
		return fmt.Errorf("неизвестный тип метрики: %s", m.MType)
	}

	if m.MType == Gauge && m.Value == nil {
		return fmt.Errorf("отсутствует значение метрики %s", m.ID)
	}

	if m.MType == Counter && m.Delta == nil && m.Value == nil {
		return fmt.Errorf("отсутствует значение метрики %s", m.ID)
	}

	return nil
}

// Normalize выполняет нормализацию полей со значениями метрики.
func (m *Metric) Normalize(val string) error {
	if m.MType == Gauge {
		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}

		m.Value = &value
		return nil
	}

	delta, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	m.Delta = &delta
	return nil
}
