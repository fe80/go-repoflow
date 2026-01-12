package repoflow

import (
	"fmt"
	"net/http"
)

// Endpoints definitions
const (
	WorkspacesEndpoint = "/1/workspaces"
)

type Workspaces struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Workspace struct {
	Id                  string `json:"id"`
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
func (c *Client) ListWorkspaces() (*[]Workspaces, error) {
	var ws []Workspaces
	err := c.DoRequest(http.MethodGet, WorkspacesEndpoint, nil, &ws)
	return &ws, err
}

// CreateWorkspace creates a new workspace with the given options
// POST /1/workspaces
func (c *Client) CreateWorkspace(opts WorkspaceOptions) (*Workspace, error) {
	var ws Workspace
	err := c.DoRequest(http.MethodPost, WorkspacesEndpoint, opts, &ws)
	return &ws, err
}

// GetWorkspace retrieves metadata for a specific workspace
// GET /1/workspaces/:id
func (c *Client) GetWorkspace(id string) (*Workspace, error) {
	var ws Workspace
	endpoint := fmt.Sprintf("%s/%s", WorkspacesEndpoint, id)
	err := c.DoRequest(http.MethodGet, endpoint, nil, &ws)
	return &ws, err
}

// DeleteWorkspace removes a workspace by its ID
// DELETE /1/workspaces/:id
func (c *Client) DeleteWorkspace(id string) (*Workspace, error) {
	var ws Workspace
	endpoint := fmt.Sprintf("%s/%s", WorkspacesEndpoint, id)
	err := c.DoRequest(http.MethodDelete, endpoint, nil, &ws)
	return &ws, err
}
