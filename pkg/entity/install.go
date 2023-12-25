package entity

import (
	"time"

	"helm.sh/helm/v3/pkg/cli/values"
)

type Parameters struct {
	Namespace          string
	Version            string
	ImageSetting       *imageSetting
	HelmChartDirectory string
	HelmOpts           values.Options
	contextName        string
	Wait               bool
	WaitDuration       time.Duration
	DryRun             bool
}

type imageSetting struct {
	Registry        string
	Repository      string
	ImagePullPolicy string
	ImagePullSecret string
	Push            bool
	ImageFileDir    string
}

type Options func(*Parameters)

func WithNamespace(namespace string) Options {
	return func(p *Parameters) {
		p.Namespace = namespace
	}
}

func WithVersion(version string) Options {
	return func(p *Parameters) {
		p.Version = version
	}
}

func WithHelmChartDirectory(helmChartDirectory string) Options {
	return func(p *Parameters) {
		p.HelmChartDirectory = helmChartDirectory
	}
}

func WithHelmOpts(helmOpts values.Options) Options {
	return func(p *Parameters) {
		p.HelmOpts = helmOpts
	}
}

func WithContextName(contextName string) Options {
	return func(p *Parameters) {
		p.contextName = contextName
	}
}

func NewParametersWithOptions(opts ...Options) *Parameters {
	p := NewParameters()
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func NewParameters() *Parameters {
	return &Parameters{
		ImageSetting: &imageSetting{},
	}
}

func (p *Parameters) IsDryRun() bool {
	return p.DryRun
}
