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
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int     `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func (m *Metric) Validate() error {
	if m.MType != Gauge && m.MType != Counter {
		return fmt.Errorf("неизвестный тип метрики: %s", m.MType)
	}

	return nil
}

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
