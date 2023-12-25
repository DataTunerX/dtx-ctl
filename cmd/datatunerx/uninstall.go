package datatunerx

import (
	"context"

	"github.com/DataTunerX/dtx-ctl/pkg/entity"
	"github.com/DataTunerX/dtx-ctl/pkg/repository/uninstall"
	"github.com/DataTunerX/dtx-ctl/pkg/valueobject"
	"github.com/spf13/cobra"
)

func addUnInstallCommand() *cobra.Command {
	unInstallParameters := entity.NewUnInstallParameters()
	unInstallChart := &cobra.Command{
		Use:   "uninstall",
		Short: "unInstall datatunerx by helm on kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			unInstallParameters.Namespace = namespace
			unInstaller := uninstall.NewUninstaller(k8sClient, unInstallParameters)
			cmd.SilenceUsage = true
			if err := unInstaller.UnInstallDataTunerXByHelm(context.Background(), k8sClient.HelmActionConfig); err != nil {
				return err
			}
			return nil
		},
	}
	unInstallChart.Flags().BoolVar(&unInstallParameters.Wait, "wait", false, "Wait for uninstallation to have completed")
	unInstallChart.Flags().DurationVar(&unInstallParameters.Timeout, "timeout", valueobject.UninstallTimeout, "Maximum time to wait for resources to be deleted")
	return unInstallChart
}
