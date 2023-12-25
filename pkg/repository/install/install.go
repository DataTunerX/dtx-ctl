package install

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DataTunerX/dtx-ctl/logging"
	"github.com/DataTunerX/dtx-ctl/pkg/entity"
	"github.com/DataTunerX/dtx-ctl/pkg/infrastructure/helm"
	"github.com/DataTunerX/dtx-ctl/pkg/infrastructure/k8s"
	"github.com/DataTunerX/dtx-ctl/pkg/valueobject"
	"github.com/DataTunerX/dtx-ctl/utils"
	"github.com/blang/semver/v4"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sInstallerImpl interface {
	GetDataTunerXStatus()
	DeleteDeployment(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error
	DeleteService(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error
	ListDeploymentByLabel(ctx context.Context, namespace string, label map[string]string) (*v1.DeploymentList, error)
	DeleteDeploymentByLabel(ctx context.Context, namespace string, label map[string]string) error
	ListServiceByLabel(ctx context.Context, namespace string, label map[string]string) (*corev1.ServiceList, error)
	DeleteValidatingWebhookConfigurationsByLabel(ctx context.Context, label map[string]string) error
	DeleteMutatingWebhookConfigurationsByLabel(ctx context.Context, label map[string]string) error
	DeleteSecretByLabel(ctx context.Context, namespace string, label map[string]string) error
	DeleteCLusterRoleBindingByLabel(ctx context.Context, label map[string]string) error
	DeleteClusterRoleByLabel(ctx context.Context, label map[string]string) error
	ListClusterRoleByLabel(ctx context.Context, label map[string]string) (*rbacv1.ClusterRoleList, error)
	DeleteServiceAccountByLabel(ctx context.Context, namespace string, label map[string]string) error
}

type Installer struct {
	client         K8sInstallerImpl
	params         *entity.Parameters
	chart          *chart.Chart
	chartVersion   semver.Version
	manifests      map[string]string
	helmYAMLValues string
}

func NewInstaller(client K8sInstallerImpl, params *entity.Parameters) (*Installer, error) {
	return &Installer{
		client: client,
		params: params,
	}, nil
}

func (i *Installer) InstallDataTunerXByHelm(ctx context.Context, k8sClient *k8s.Client) error {
	userSetValues, err := i.params.HelmOpts.MergeValues(getter.All(cli.New()))
	if err != nil {
		return err
	}
	logging.Logger.Info("get user set values", slog.Any("userSetValues:", userSetValues))
	version, err := utils.ParseDataTunerXVersion(i.params.Version)
	if err != nil {
		return err
	}
	i.chartVersion = version
	switch {
	case helm.CheckVersion("=0.0.1")(i.chartVersion):
		logging.Logger.Info("datatunerx version check", slog.String("version", i.chartVersion.String()))
		if i.params.ImageSetting != nil {
			if i.params.ImageSetting.Registry != "" {
				userSetValues["global.registry"] = i.params.ImageSetting.Registry
			}
			if i.params.ImageSetting.Repository != "" {
				userSetValues["global.repository"] = i.params.ImageSetting.Repository
			}

			if i.params.ImageSetting.ImagePullPolicy != "" {
				userSetValues["global.imagePullPolicy"] = i.params.ImageSetting.ImagePullPolicy
			}

			if i.params.ImageSetting.ImagePullSecret != "" {
				userSetValues["global.imagePullSecret"] = i.params.ImageSetting.ImagePullSecret
			}

			if i.params.ImageSetting.Push {
				// todo(tigerK) push  client
			}

			if i.params.ImageSetting.ImageFileDir != "" {
				// todo(tigerK) load image client
			}
		}
		datatunerxChart, err := helm.GetChartByVersion(version)
		if err != nil {
			return err
		}
		i.chart = datatunerxChart
	default:
		return fmt.Errorf("datatunerx version unsupported %s", i.chartVersion.String())
	}
	helmClient := action.NewInstall(k8sClient.HelmActionConfig)
	helmClient.ReleaseName = valueobject.HelmReleaseName
	logging.Logger.Info("get", slog.String("releaseName", valueobject.HelmReleaseName))
	helmClient.Namespace = i.params.Namespace
	logging.Logger.Info("set install", slog.String("namespace", helmClient.Namespace))
	helmClient.Wait = i.params.Wait
	helmClient.Timeout = i.params.WaitDuration
	helmClient.DryRun = i.params.IsDryRun()
	logging.Logger.Info("start run helm install")
	release, err := helmClient.RunWithContext(ctx, i.chart, userSetValues)
	if err != nil {
		return err
	}
	logging.Logger.Info("install success", slog.String("releaseName", release.Name))
	return nil
}
