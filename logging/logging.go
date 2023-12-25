package logging

import (
	"log/slog"

	"github.com/pterm/pterm"
)

var Logger *slog.Logger

func init() {
	handler := pterm.NewSlogHandler(&pterm.DefaultLogger)
	Logger = slog.New(handler)
}
