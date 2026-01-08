package factory

import (
	"log/slog"
	"repoflow/pkg/config"
)

type Utils struct {
	Cfg    *config.Config
	Logger *slog.Logger
	Output string
}
