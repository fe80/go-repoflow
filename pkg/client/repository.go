package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Endpoints definitions
const (
	RepositoryEndpoint = "/repositories"
)

type Repositories struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	PackageType    string `json:"packageType"`
	RepositoryType string `json:"repositoryType"`
	Status         string `json:"status"`
}

type UploadTargetLocalRepository struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ChildRepository struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Repository struct {
	RepositoryType                    string                      `json:"repositoryType"`
	Status                            string                      `json:"status"`
	WorkspaceId                       string                      `json:"workspaceId"`
	UploadTargetLocalRepository       UploadTargetLocalRepository `json:"uploadTargetLocalRepository"`
	ChildRepositories                 []ChildRepository           `json:"childRepositories"`
	RemoteRepositoryUrl               *string                     `json:"remoteRepositoryUrl,omitempty"`
	IsRemoteCacheEnabled              *string                     `json:"isRemoteCacheEnabled,omitempty"`
	FileCacheTimeTillRevalidation     *int                        `json:"fileCacheTimeTillRevalidation,omitempty"`
	MetadataCacheTimeTillRevalidation *int                        `json:"metadataCacheTimeTillRevalidation,omitempty"`
}

type RepositoryOptions struct {
	Name        string `json:"name"`
	PackageType string `json:"packageType"`
}

// RepositoryRemoteRemote defines the payload for creating a remote repository
type RepositoryRemoteOptions struct {
	Name                              string `json:"name"`
	PackageType                       string `json:"packageType"`
	RemoteRepositoryUrl               string `json:"remoteRepositoryUrl"`
	RemoteRepositoryUsername          string `json:"remoteRepositoryUsername"`
	RemoteRepositoryPassword          string `json:"remoteRepositoryPassword"`
	IsRemoteCacheEnabled              bool   `json:"isRemoteCacheEnabled"`
	FileCacheTimeTillRevalidation     *int   `json:"fileCacheTimeTillRevalidation"`
	MetadataCacheTimeTillRevalidation *int   `json:"metadataCacheTimeTillRevalidation"`
}

// RepositoryVirtualRemote defines the payload for creating a virtual repository
type RepositoryVirtualOptions struct {
	Name                    string   `json:"name"`
	PackageType             string   `json:"packageType"`
	ChildRepositoryIds      []string `json:"childRepositoryIds"`
	UploadLocalRepositoryId *string  `json:"uploadLocalRepositoryId,omitempty"`
}

type RepostotryDelete struct {
	RepositoryId string `json:"repositoryId"`
	Status       string `json:"status"`
}

// ListRepositories retrieves all available repository
// GET /1/workspaces/:workspace/repositories
func (c *Client) ListRepositories(workspace string) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	u := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// GetRepository retrieves metadata for a specific repository
// GET /1/workspaces/:workspace/repositories/:id
func (c *Client) GetRepository(workspace string, id string) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	u := fmt.Sprintf("%s%s/%s", c.BaseURL, endpoint, id)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// ListRepositoryPackages list all available package in a repository
// GET /1/workspaces/:workspace/repositories/:id/packages
func (c *Client) ListRepositoryPackages(workspace string, id string) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	u := fmt.Sprintf("%s%s/%s/packages", c.BaseURL, endpoint, id)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// CreateRepository create a new repository with the given options
// POST /1/workspaces/:workspace/repositories/:store
func (c *Client) CreateRepository(workspace string, store string, opts any) (*http.Response, error) {
	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	u := fmt.Sprintf("%s%s/%s", c.BaseURL, endpoint, store)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// DeleteRepository removes a workspace by its ID
// DELETE /1/workspaces/:id/repositories/:id
func (c *Client) DeleteRepository(workspace string, id string) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	u := fmt.Sprintf("%s%s/%s", c.BaseURL, endpoint, id)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}

// DeleteRepositoryContent removes a workspace by its ID
// DELETE /1/workspaces/:id/repositories/:id
func (c *Client) DeleteRepositoryContent(workspace string, id string) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	u := fmt.Sprintf("%s%s/%s/content", c.BaseURL, endpoint, id)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.HTTPClient.Do(req)
}
