package selectel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ServersClient представляет клиент для работы с API выделенных серверов Selectel
type ServersClient struct {
	HTTPClient *http.Client
	Token      string
	BaseURL    string
	UserAgent  string
	ctx        context.Context
}

// ServersClientOptions содержит опции для создания клиента серверов
type ServersClientOptions struct {
	Token      string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Context    context.Context
}

// NewServersClient создает новый экземпляр клиента для работы с выделенными серверами
func NewServersClient(options *ServersClientOptions) (*ServersClient, error) {
	if options == nil {
		return nil, fmt.Errorf("client options cannot be nil")
	}

	if options.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	client := &ServersClient{
		Token: options.Token,
	}

	// Устанавливаем базовый URL
	if options.BaseURL != "" {
		client.BaseURL = options.BaseURL
	} else {
		client.BaseURL = "https://api.selectel.ru/servers/v2/"
	}

	// Устанавливаем User-Agent
	if options.UserAgent != "" {
		client.UserAgent = options.UserAgent
	} else {
		client.UserAgent = "terraform-provider-selectel/1.0"
	}

	// Устанавливаем HTTP клиент
	if options.HTTPClient != nil {
		client.HTTPClient = options.HTTPClient
	} else {
		client.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// Устанавливаем контекст
	if options.Context != nil {
		client.ctx = options.Context
	} else {
		client.ctx = context.Background()
	}

	return client, nil
}

// DoRequest выполняет HTTP запрос к API выделенных серверов
func (c *ServersClient) DoRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = u.Path + path

	log.Printf("[DEBUG] Making %s request to: %s", method, u.String())
	log.Printf("[DEBUG] Using token (length: %d): %s...", len(c.Token), c.Token[:min(20, len(c.Token))])
	log.Printf("[DEBUG] FULL TOKEN FOR DEBUGGING: %s", c.Token)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		log.Printf("[DEBUG] Request body JSON: %s", string(jsonBody))
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	log.Printf("[DEBUG] Response status: %d", resp.StatusCode)

	return resp, nil
}

// ParseResponse парсит ответ от API
func (c *ServersClient) ParseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError ServersAPIError
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return &apiError
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// ServersAPIError представляет ошибку API выделенных серверов
type ServersAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error реализует интерфейс error
func (e *ServersAPIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// ServersListOptions содержит опции для запросов списка серверов
type ServersListOptions struct {
	Page     int    `url:"page,omitempty"`
	PerPage  int    `url:"per_page,omitempty"`
	Sort     string `url:"sort,omitempty"`
	Status   string `url:"status,omitempty"`
	Location string `url:"location,omitempty"`
}

// BuildQueryString строит строку запроса из опций
func (opts *ServersListOptions) BuildQueryString() string {
	values := url.Values{}

	if opts.Page > 0 {
		values.Add("page", strconv.Itoa(opts.Page))
	}

	if opts.PerPage > 0 {
		values.Add("per_page", strconv.Itoa(opts.PerPage))
	}

	if opts.Sort != "" {
		values.Add("sort", opts.Sort)
	}

	if opts.Status != "" {
		values.Add("status", opts.Status)
	}

	if opts.Location != "" {
		values.Add("location", opts.Location)
	}

	if len(values) > 0 {
		return "?" + values.Encode()
	}

	return ""
}

// min возвращает минимальное из двух чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
