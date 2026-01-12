package factory

import (
	"github.com/fe80/go-repoflow/pkg/repoflow"
	"github.com/fe80/go-repoflow/pkg/config"
)

func GetClient(cfg *config.Config) *repoflow.Client {
	return repoflow.NewClient(cfg.URL, cfg.Token)
}
