package http

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/iamamatkazin/metrics.git/pkg/config"
)

type Client struct {
	*http.Client
	cfg *config.Config
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		Client: &http.Client{
			Timeout:   cfg.Client.Timeout,
			Transport: &http.Transport{},
		},
	}
}

func (c *Client) Post(ctx context.Context, url, contentType string) error {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Client.Timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, http.NoBody)
	if err != nil {
		return err
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", contentType)

	// отправляем запрос и получаем ответ
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
