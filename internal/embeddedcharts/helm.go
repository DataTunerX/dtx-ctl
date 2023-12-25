package embeddedcharts

import (
	"embed"
)

var (
	//go:embed *.tgz
	EmbeddedChartsFS embed.FS
)
