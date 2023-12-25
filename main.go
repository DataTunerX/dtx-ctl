package main

import (
	"log/slog"
	"os"

	"github.com/DataTunerX/dtx-ctl/cmd/datatunerx"
	"github.com/DataTunerX/dtx-ctl/logging"
)

func main() {
	if err := datatunerx.NewDataTunerXCommand().Execute(); err != nil {
		logging.Logger.Error("An error occurred", slog.String("Err", err.Error()))
		os.Exit(1)
	}
}
