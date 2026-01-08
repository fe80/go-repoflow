package factory

import (
	"github.com/fe80/go-repoflow/pkg/client"
	"github.com/fe80/go-repoflow/pkg/config"
)

func GetClient(cfg *config.Config) *client.Client {
	return client.NewClient(cfg.URL, cfg.Token)
}
