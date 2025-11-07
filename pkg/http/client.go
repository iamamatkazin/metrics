package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iamamatkazin/metrics.git/pkg/config/agent"
)

type Client struct {
	*http.Client
	cfg *agent.Config
}

func New(cfg *agent.Config) *Client {
	return &Client{
		cfg: cfg,
		Client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: &http.Transport{},
		},
	}
}

func (c *Client) Post(ctx context.Context, url, contentType string, data any) (err error) {
	var (
		request *http.Request
		body    []byte
	)

	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	if data == nil {
		request, err = http.NewRequestWithContext(ctx, http.MethodPost, url, http.NoBody)
	} else {
		body, err = json.Marshal(data)
		if err != nil {
			return err
		}
		request, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	}
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
