package factory

import (
	"github.com/fe80/go-repoflow/pkg/repoflow"
	"github.com/fe80/go-repoflow/pkg/config"
	"log/slog"
)

type Utils struct {
	Cfg       *config.Config
	Logger    *slog.Logger
	Output    string
	apiClient *repoflow.Client
}

func (u *Utils) GetAPIClient() *repoflow.Client {
	if u.apiClient == nil {
		u.apiClient = GetClient(u.Cfg)
	}
	return u.apiClient
}
