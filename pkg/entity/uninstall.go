package entity

import (
	"io"
	"time"
)

type UnInstallParameters struct {
	Namespace            string
	Version              string
	Writer               io.Writer
	Wait                 bool
	HelmValuesSecretName string
	RedactHelmCertKeys   bool
	HelmChartDirectory   string
	WorkerCount          int
	Timeout              time.Duration
}

func NewUnInstallParameters() *UnInstallParameters {
	return &UnInstallParameters{}
}
