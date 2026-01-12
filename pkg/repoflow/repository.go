package repoflow

import (
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

type PackageRepository struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ChildRepository struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Repository struct {
	Name                              string                      `json:"name"`
	Id                                string                      `json:"id"`
	RepositoryType                    string                      `json:"repositoryType"`
	PackageType                       string                      `json:"packageType"`
	Status                            string                      `json:"status"`
	WorkspaceId                       string                      `json:"workspaceId"`
	UploadTargetLocalRepository       UploadTargetLocalRepository `json:"uploadTargetLocalRepository"`
	ChildRepositories                 []ChildRepository           `json:"childRepositories"`
	RemoteRepositoryUrl               *string                     `json:"remoteRepositoryUrl,omitempty"`
	IsRemoteCacheEnabled              bool                        `json:"isRemoteCacheEnabled"`
	FileCacheTimeTillRevalidation     *int                        `json:"fileCacheTimeTillRevalidation,omitempty"`
	MetadataCacheTimeTillRevalidation *int                        `json:"metadataCacheTimeTillRevalidation,omitempty"`
}

type RepositoryOptions struct {
	Name        string `json:"name"`
	PackageType string `json:"packageType"`
}

type RepositoryPackages struct {
	Total    int                  `json:"total"`
	Offset   int                  `json:"offset"`
	Limit    int                  `json:"int"`
	Packages []*PackageRepository `json:"packages"`
}

// RepositoryRemoteRemote defines the payload for creating a remote repository
type RepositoryRemoteOptions struct {
	Name                              string `json:"name"`
	PackageType                       string `json:"packageType"`
	RemoteRepositoryUrl               string `json:"remoteRepositoryUrl"`
	IsRemoteCacheEnabled              bool   `json:"isRemoteCacheEnabled"`
	RemoteRepositoryUsername          string `json:"remoteRepositoryUsername,omitempty"`
	RemoteRepositoryPassword          string `json:"remoteRepositoryPassword,omitempty"`
	FileCacheTimeTillRevalidation     *int   `json:"fileCacheTimeTillRevalidation,omitempty"`
	MetadataCacheTimeTillRevalidation *int   `json:"metadataCacheTimeTillRevalidation,omitempty"`
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
func (c *Client) ListRepositories(workspace string) (*[]Repositories, error) {
	var rep []Repositories
	endpoint := fmt.Sprintf("%s/%s%s", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	err := c.DoRequest(http.MethodGet, endpoint, nil, &rep)
	return &rep, err
}

// GetRepository retrieves metadata for a specific repository
// GET /1/workspaces/:workspace/repositories/:id
func (c *Client) GetRepository(workspace string, id string) (*Repository, error) {
	var rep Repository
	endpoint := fmt.Sprintf("%s/%s%s/%s", WorkspacesEndpoint, workspace, RepositoryEndpoint, id)
	err := c.DoRequest(http.MethodGet, endpoint, nil, &rep)
	return &rep, err
}

// ListRepositoryPackages list all available package in a repository
// GET /1/workspaces/:workspace/repositories/:id/packages
func (c *Client) ListRepositoryPackages(workspace string, id string) (*RepositoryPackages, error) {
	var rep RepositoryPackages
	endpoint := fmt.Sprintf("%s/%s%s/%s/packages", WorkspacesEndpoint, workspace, RepositoryEndpoint, id)
	err := c.DoRequest(http.MethodGet, endpoint, nil, &rep)
	return &rep, err
}

// CreateRepository create a new repository with the given options
// POST /1/workspaces/:workspace/repositories/:store
func (c *Client) CreateRepository(workspace string, store string, opts any) (*Repositories, error) {
	var rep Repositories
	endpoint := fmt.Sprintf("%s/%s%s/%s", WorkspacesEndpoint, workspace, RepositoryEndpoint, store)
	err := c.DoRequest(http.MethodPost, endpoint, opts, &rep)
	return &rep, err
}

// CreateLocalRepository create a new repository with the given options
// POST /1/workspaces/:workspace/repositories/local
func (c *Client) CreateLocalRepository(workspace string, opts RepositoryOptions) (*Repositories, error) {
	var rep Repositories
	endpoint := fmt.Sprintf("%s/%s%s/local", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	err := c.DoRequest(http.MethodPost, endpoint, opts, &rep)
	return &rep, err
}

// CreateRemoteRepository create a new repository with the given options
// POST /1/workspaces/:workspace/repositories/remote
func (c *Client) CreateRemoteRepository(workspace string, opts RepositoryRemoteOptions) (*Repositories, error) {
	var rep Repositories
	endpoint := fmt.Sprintf("%s/%s%s/remote", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	err := c.DoRequest(http.MethodPost, endpoint, opts, &rep)
	return &rep, err
}

// CreateVirtualRepository create a new repository with the given options
// POST /1/workspaces/:workspace/repositories/virtual
func (c *Client) CreateVirtualRepository(workspace string, opts RepositoryVirtualOptions) (*Repositories, error) {
	var rep Repositories
	endpoint := fmt.Sprintf("%s/%s%s/virtual", WorkspacesEndpoint, workspace, RepositoryEndpoint)
	err := c.DoRequest(http.MethodPost, endpoint, opts, &rep)
	return &rep, err
}

// DeleteRepository removes a workspace by its ID
// DELETE /1/workspaces/:id/repositories/:id
func (c *Client) DeleteRepository(workspace string, id string) (*RepostotryDelete, error) {
	var rep RepostotryDelete
	endpoint := fmt.Sprintf("%s/%s%s/%s", WorkspacesEndpoint, workspace, RepositoryEndpoint, id)
	err := c.DoRequest(http.MethodDelete, endpoint, nil, &rep)
	return &rep, err
}

// DeleteRepositoryContent removes a workspace by its ID
// DELETE /1/workspaces/:id/repositories/:id
func (c *Client) DeleteRepositoryContent(workspace string, id string) (*RepostotryDelete, error) {
	var rep RepostotryDelete
	endpoint := fmt.Sprintf("%s/%s%s/%s/content", WorkspacesEndpoint, workspace, RepositoryEndpoint, id)
	err := c.DoRequest(http.MethodDelete, endpoint, nil, &rep)
	return &rep, err
}
