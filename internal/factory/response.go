package factory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fe80/go-repoflow/pkg/client"
)

// handleResponse processes the API response.
// If it's an error, it returns an error with the body content.
// If it's a success, it returns nil so the caller can proceed with decoding.
func (u *Utils) HandleResponse(resp *http.Response) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("API Error (%s): failed to read error body", resp.Status)
	}

	u.Logger.Debug("API Response received",
		"method", resp.Request.Method,
		"url", resp.Request.URL.String(),
		"status", resp.Status,
		"payload", string(body),
	)

	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if resp.StatusCode < 400 {
		return nil
	}

	var apiErr client.APIError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
		return fmt.Errorf("API Error (%d): %s - %s", resp.StatusCode, apiErr.Code, apiErr.Message)
	}

	var apiErrs client.APIErrors
	if err := json.Unmarshal(body, &apiErrs); err == nil && len(apiErrs.Errors) > 0 {
		return fmt.Errorf("API Errors (%d): %s", resp.StatusCode, strings.Join(apiErrs.Errors, "; "))
	}

	// Returning the error stops the command execution and displays the message
	return fmt.Errorf("request failed with status: %s", resp.Status)
}
