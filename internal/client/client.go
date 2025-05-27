package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/models"
)

const (
	// DefaultEndpoint - базовый URL API Selectel для выделенных серверов
	DefaultEndpoint = "https://api.selectel.ru/dedicated/v2"

	// DefaultTimeout - таймаут по умолчанию для HTTP запросов
	DefaultTimeout = 30 * time.Second

	// DefaultRetryAttempts - количество попыток повтора запроса
	DefaultRetryAttempts = 3

	// DefaultRetryDelay - задержка между попытками
	DefaultRetryDelay = 1 * time.Second
)

// Client представляет HTTP клиент для API Selectel
type Client struct {
	httpClient  *http.Client
	endpoint    string
	token       string
	projectID   string
	userAgent   string
	retryConfig RetryConfig
}

// RetryConfig содержит настройки повторных попыток
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64 // множитель для экспоненциального backoff
}

// Config содержит конфигурацию клиента
type Config struct {
	Endpoint    string
	Token       string
	ProjectID   string
	Timeout     time.Duration
	RetryConfig *RetryConfig
	UserAgent   string
}

// NewClient создает новый HTTP клиент
func NewClient(config Config) *Client {
	if config.Endpoint == "" {
		config.Endpoint = DefaultEndpoint
	}

	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	if config.RetryConfig == nil {
		config.RetryConfig = &RetryConfig{
			MaxAttempts: DefaultRetryAttempts,
			Delay:       DefaultRetryDelay,
			Backoff:     2.0,
		}
	}

	if config.UserAgent == "" {
		config.UserAgent = "terraform-provider-selectel-baremetal/dev"
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		httpClient:  httpClient,
		endpoint:    strings.TrimSuffix(config.Endpoint, "/"),
		token:       config.Token,
		projectID:   config.ProjectID,
		userAgent:   config.UserAgent,
		retryConfig: *config.RetryConfig,
	}
}

// buildURL создает полный URL для запроса
func (c *Client) buildURL(path string) string {
	return c.endpoint + "/" + strings.TrimPrefix(path, "/")
}

// prepareRequest подготавливает HTTP запрос
func (c *Client) prepareRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.buildURL(path)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+c.token)

	if c.projectID != "" {
		req.Header.Set("X-Project-ID", c.projectID)
	}

	return req, nil
}

// executeRequest выполняет HTTP запрос с повторными попытками
func (c *Client) executeRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	delay := c.retryConfig.Delay

	for attempt := 0; attempt < c.retryConfig.MaxAttempts; attempt++ {
		if attempt > 0 {
			tflog.Debug(ctx, "Retrying request", map[string]interface{}{
				"attempt": attempt + 1,
				"delay":   delay,
				"url":     req.URL.String(),
			})

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
				delay = time.Duration(float64(delay) * c.retryConfig.Backoff)
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			tflog.Warn(ctx, "Request failed", map[string]interface{}{
				"error":   err.Error(),
				"attempt": attempt + 1,
				"url":     req.URL.String(),
			})
			continue
		}

		// Проверяем, нужно ли повторить запрос
		if c.shouldRetry(resp.StatusCode) {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			tflog.Warn(ctx, "Request returned retryable status", map[string]interface{}{
				"status_code": resp.StatusCode,
				"attempt":     attempt + 1,
				"url":         req.URL.String(),
			})
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryConfig.MaxAttempts, lastErr)
}

// shouldRetry определяет, нужно ли повторить запрос на основе статус кода
func (c *Client) shouldRetry(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	default:
		return false
	}
}

// parseResponse парсит ответ API
func (c *Client) parseResponse(ctx context.Context, resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	tflog.Debug(ctx, "API Response", map[string]interface{}{
		"status_code": resp.StatusCode,
		"body_size":   len(body),
		"url":         resp.Request.URL.String(),
	})

	// Проверяем статус код
	if resp.StatusCode >= 400 {
		var apiError models.APIError
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return &apiError
	}

	// Парсим успешный ответ
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// DoRequest выполняет HTTP запрос
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	req, err := c.prepareRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Making API request", map[string]interface{}{
		"method": method,
		"url":    req.URL.String(),
	})

	resp, err := c.executeRequest(ctx, req)
	if err != nil {
		return err
	}

	return c.parseResponse(ctx, resp, result)
}

// Get выполняет GET запрос
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.DoRequest(ctx, http.MethodGet, path, nil, result)
}

// Post выполняет POST запрос
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.DoRequest(ctx, http.MethodPost, path, body, result)
}

// Put выполняет PUT запрос
func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.DoRequest(ctx, http.MethodPut, path, body, result)
}

// Patch выполняет PATCH запрос
func (c *Client) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.DoRequest(ctx, http.MethodPatch, path, body, result)
}

// Delete выполняет DELETE запрос
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.DoRequest(ctx, http.MethodDelete, path, nil, nil)
}

// GetWithQuery выполняет GET запрос с query параметрами
func (c *Client) GetWithQuery(ctx context.Context, path string, params url.Values, result interface{}) error {
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	return c.Get(ctx, path, result)
}
