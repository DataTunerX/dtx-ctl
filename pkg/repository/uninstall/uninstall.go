package uninstall

import (
	"context"
	"log/slog"

	"github.com/DataTunerX/dtx-ctl/logging"
	"github.com/DataTunerX/dtx-ctl/pkg/entity"
	"github.com/DataTunerX/dtx-ctl/pkg/repository/install"
	"github.com/DataTunerX/dtx-ctl/pkg/valueobject"
	"github.com/blang/semver/v4"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Uninstaller struct {
	client  install.K8sInstallerImpl
	params  *entity.UnInstallParameters
	version semver.Version
}

func NewUninstaller(client install.K8sInstallerImpl, p *entity.UnInstallParameters) *Uninstaller {
	return &Uninstaller{
		client: client,
		params: p,
	}
}

func (u *Uninstaller) UnInstallDataTunerXByHelm(ctx context.Context, actionConfig *action.Configuration) error {
	logging.Logger.Info("start uninstall", slog.String("releaseName", valueobject.HelmReleaseName))
	helmClient := action.NewUninstall(actionConfig)
	helmClient.Wait = u.params.Wait
	helmClient.Timeout = u.params.Timeout
	if _, err := helmClient.Run(valueobject.HelmReleaseName); err != nil {
		return err
	}
	logging.Logger.Info("start delete deployment")
	if err := u.client.DeleteDeploymentByLabel(ctx, u.params.Namespace, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("start delete service")
	serviceList, err := u.client.ListServiceByLabel(ctx, u.params.Namespace, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue})
	if err != nil {
		return err
	}
	for i := range serviceList.Items {
		service := serviceList.Items[i]
		if err := u.client.DeleteService(ctx, service.Name, u.params.Namespace, metav1.DeleteOptions{}); err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return err
		}
	}
	logging.Logger.Info("start delete secret")
	if err := u.client.DeleteSecretByLabel(ctx, u.params.Namespace, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("start delete clusterRoleBinding")
	if err := u.client.DeleteCLusterRoleBindingByLabel(ctx, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("start delete clusterRole")
	if err := u.client.DeleteClusterRoleByLabel(ctx, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("start delete mutatingWebhookConfiguration")
	if err := u.client.DeleteMutatingWebhookConfigurationsByLabel(ctx, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("start delete validatingWebhookConfiguration")
	if err := u.client.DeleteValidatingWebhookConfigurationsByLabel(ctx, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("start delete serviceAccount")
	if err := u.client.DeleteServiceAccountByLabel(ctx, u.params.Namespace, map[string]string{valueobject.ComposeCommonLabelKey: valueobject.ComposeCommonLabelValue}); err != nil {
		return err
	}
	logging.Logger.Info("success uninstall", slog.String("releaseName", valueobject.HelmReleaseName))
	return nil
}
