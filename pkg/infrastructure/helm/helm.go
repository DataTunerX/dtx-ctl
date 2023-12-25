package helm

import (
	"bytes"
	"fmt"

	"github.com/DataTunerX/dtx-ctl/internal/embeddedcharts"
	"github.com/blang/semver/v4"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func GetChartByVersion(version semver.Version) (*chart.Chart, error) {
	datatunerxChart, err := getChartByEmbedded(version)
	if err != nil {
		return nil, err
	}
	// todo(tigerK) If it doesn't exist locally, pull it from the repository
	return datatunerxChart, nil
}

func getChartByEmbedded(version semver.Version) (*chart.Chart, error) {
	chartTgz, err := embeddedcharts.EmbeddedChartsFS.ReadFile(fmt.Sprintf("datatunerx-%s.tgz", version))
	if err != nil {
		return nil, err
	}
	return loader.LoadArchive(bytes.NewReader(chartTgz))
}

func CheckVersion(version string) semver.Range {
	ver, _ := semver.ParseRange(version)
	return ver
}
