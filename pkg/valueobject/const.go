package valueobject

import "time"

const (
	StatusWaitDuration = 10 * time.Second
	UninstallTimeout   = 5 * time.Minute
	InstallTimeout     = 5 * time.Minute
	HelmReleaseName    = "datatunerx"
)

const (
	ComposeCommonLabelKey   = "datatunerx.io/part-of"
	ComposeCommonLabelValue = "datatunerx"
)
