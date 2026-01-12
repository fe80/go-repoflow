package repoflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIErrors struct {
	Errors []string `json:"errors"`
}

func NewClient(baseUrl string, token string) *Client {
	return &Client{
		BaseURL: baseUrl,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
}

func (e *APIErrors) Error() string {
	if len(e.Errors) == 0 {
		return "unknown api error"
	}
	return fmt.Sprintf("api error: %v", strings.Join(e.Errors, "; "))
}

func (c *Client) DoRequest(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errs APIErrors
		bodyBytes, _ := io.ReadAll(resp.Body)

		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &errs); err == nil && len(errs.Errors) > 0 {
				return &errs
			}
		}

		return fmt.Errorf("api error: status %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
