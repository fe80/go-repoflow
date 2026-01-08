package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Endpoints definitions
const (
	WorkspacesEndpoint = "/1/workspaces"
)

type Workspaces struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Workspace struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	StorageUsageInByte  int    `json:"storageUsageInByte"`
	StorageLimitInByte  *int   `json:"storageLimitInByte"`
	TransferUsageInByte int    `json:"transferUsageInByte"`
	TransferLimitInByte *int   `json:"transferLimitInByte"`
	PackageUsage        int    `json:"packageUsage"`
	PackageLimit        *int   `json:"packageLimit"`
	AiUsageCount        int    `json:"aiUsageCount"`
	AiUsageLimit        *int   `json:"aiUsageLimit"`
}

// WorkspaceOptions defines the payload for creating a workspace
type WorkspaceOptions struct {
	Name           string  `json:"name"`
	PackageLimit   *int    `json:"packageLimit,omitempty"`
	BandwidthLimit *int    `json:"bandwidthLimit,omitempty"`
	StorageLimit   *int    `json:"storageLimit,omitempty"`
	Comments       *string `json:"comment,omitempty"`
}

// ListWorkspaces retrieves all available workspaces
// GET /1/workspaces
func (c *Client) ListWorkspaces() (*http.Response, error) {
	u := fmt.Sprintf("%s%s", c.BaseURL, WorkspacesEndpoint)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// CreateWorkspace creates a new workspace with the given options
// POST /1/workspaces
func (c *Client) CreateWorkspace(opts WorkspaceOptions) (*http.Response, error) {
	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s%s", c.BaseURL, WorkspacesEndpoint)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// GetWorkspace retrieves metadata for a specific workspace
// GET /1/workspaces/:id
func (c *Client) GetWorkspace(id string) (*http.Response, error) {
	u := fmt.Sprintf("%s%s/%s", c.BaseURL, WorkspacesEndpoint, id)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// DeleteWorkspace removes a workspace by its ID
// DELETE /1/workspaces/:id
func (c *Client) DeleteWorkspace(id string) (*http.Response, error) {
	u := fmt.Sprintf("%s%s/%s", c.BaseURL, WorkspacesEndpoint, id)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}
