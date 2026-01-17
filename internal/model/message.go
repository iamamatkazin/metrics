package model

type Message struct {
	Date    int64    `json:"ts"`
	Metrics []string `json:"metrics"`
	IP      string   `json:"ip_address"`
}
