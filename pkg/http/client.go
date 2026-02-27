// Package http предоставляет HTTP клиент для системы метрик.
// Реализует кастомный HTTP клиент с логикой повторных попыток и HMAC подписью
// для безопасной передачи метрик.
package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/common"
)

// Clienter - интерфейс отправки HTTP запросов.
type Clienter interface {
	Post(ctx context.Context, url, contentType, key string, data []byte) (err error)
}

// Client - структура реализации кастомного HTTP клиента.
type Client struct {
	*http.Client
	ip string
}

// New создает новый экземпляр Client с заданным таймаутом.
func New(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout:   timeout,
			Transport: &http.Transport{},
		},
		ip: getIPAdress(),
	}
}

// Post отправляет HTTP POST запрос с метриками на сервер.
// Поддерживает повторные попытки при временных сбоях и автоматически
// добавляет HMAC-SHA256 подпись если указан ключ.
func (c *Client) Post(ctx context.Context, url, contentType, key string, data []byte) (err error) {
	var (
		request *http.Request
	)

	if data == nil {
		request, err = newRequestWithContext(ctx, http.MethodPost, url, http.NoBody)
	} else {
		request, err = newRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	}
	if err != nil {
		return err
	}

	// В заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("X-REAL-IP", c.ip)

	if key != "" {
		request.Header.Set("HashSHA256", common.CalcSign([]byte(key), data))
	}

	// Отправляем запрос и получаем ответ
	response, err := c.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка обработки запроса с кодом: %d", response.StatusCode)
	}

	return nil
}

// newRequestWithContext создает HTTP запрос с контекстом и реализует механизм повторных попыток.
func newRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	timerRetriable := time.NewTimer(0)
	count := 0

	for {
		select {
		case <-ctx.Done():
		case <-timerRetriable.C:
			request, err := http.NewRequestWithContext(ctx, method, url, body)
			if err != nil {
				if count > 3 {
					return nil, err
				}

				timerRetriable.Reset(time.Duration(2*count+1) * time.Second)
				count++
			}

			return request, nil
		}
	}
}

func getIPAdress() string {
	host, err := os.Hostname()
	if err != nil {
		slog.Error("Ошибка получения имени хоста:", slog.Any("error", err))
		return ""
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		slog.Error("Ошибка получения ip адреса:", slog.Any("error", err))
		return ""
	}

	for _, addr := range addrs {
		return addr.To4().String()
	}

	return ""
}
