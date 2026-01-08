package factory

import (
	"repoflow/pkg/client"
	"repoflow/pkg/config"
)

func GetClient(cfg *config.Config) *client.Client {
	return client.NewClient(cfg.URL, cfg.Token)
}
