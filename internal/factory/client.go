package factory

import (
	"github.com/fe80/go-repoflow/pkg/config"
	"github.com/fe80/go-repoflow/pkg/repoflow"
)

func GetClient(cfg *config.Config) *repoflow.Client {
	return repoflow.NewClient(cfg.URL, cfg.Token)
}
