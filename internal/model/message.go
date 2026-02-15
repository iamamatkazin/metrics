// Package model предоставляет структуры данных для системы метрик.
// Содержит типы Metric и Message для представления метрик и
// сообщений аудита соответственно.
package model

// Message представляет сообщение аудита, содержащее информацию о
// том, какие метрики были изменены и с какого IP адреса.
// generate:reset
type Message struct {
	IP      string   `json:"ip_address"`
	Metrics []string `json:"metrics"`
	Date    int64    `json:"ts"`
}
