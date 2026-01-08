package factory

import (
	"log/slog"
	"repoflow/pkg/config"
	"repoflow/pkg/client"
)

type Utils struct {
	Cfg    *config.Config
	Logger *slog.Logger
	Output string
	apiClient *client.Client
}

func (u *Utils) GetAPIClient() *client.Client {
    if u.apiClient == nil {
        u.apiClient = GetClient(u.Cfg)
    }
    return u.apiClient
}
