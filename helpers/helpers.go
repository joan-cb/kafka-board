package helpers

import (
	"log/slog"
)

type Helpers struct {
	logger *slog.Logger
}

func ReturnHelpers(logger *slog.Logger) *Helpers {
	return &Helpers{logger: logger}
}
