package datatunerx

import (
	"context"

	"github.com/DataTunerX/dtx-ctl/pkg/entity"
	"github.com/DataTunerX/dtx-ctl/pkg/repository/install"
	"github.com/DataTunerX/dtx-ctl/pkg/valueobject"
	"github.com/spf13/cobra"
)

var installChart *cobra.Command

func addInstallCommand() *cobra.Command {
	params := entity.NewParametersWithOptions()
	installChart = &cobra.Command{
		Use:   "install",
		Short: "Install datatunerx by helm on kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			params.Namespace = namespace
			installer, err := install.NewInstaller(k8sClient, params)
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true
			if err := installer.InstallDataTunerXByHelm(context.Background(), k8sClient); err != nil {
				return err
			}
			return nil
		},
	}
	installChart.Flags().StringVar(&params.HelmChartDirectory, "chart-directory", "", "Helm chart directory")
	installChart.Flags().StringSliceVarP(&params.HelmOpts.ValueFiles, "values", "f", []string{}, "Specify helm values in a YAML file or a URL (can specify multiple)")
	installChart.Flags().StringArrayVar(&params.HelmOpts.Values, "set", []string{}, "Set helm values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	installChart.Flags().StringArrayVar(&params.HelmOpts.StringValues, "set-string", []string{}, "Set helm STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	installChart.Flags().StringArrayVar(&params.HelmOpts.FileValues, "set-file", []string{}, "Set helm values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")
	installChart.Flags().StringVar(&params.Version, "version", "0.0.1", "Specify a version constraint for the chart version to use. This constraint can be anything that the Helm client supports, for example --version=0.0.1 ")
	installChart.Flags().StringVar(&params.ImageSetting.Registry, "registry", "", "Specify a registry for the chart version to use. This constraint can be anything that the Helm client supports, for example --registry=registry.cn-hangzhou.aliyuncs.com ")
	installChart.Flags().StringVar(&params.ImageSetting.Repository, "repository", "", "Specify a repository for the chart version to use. This constraint can be anything that the Helm client supports, for example --repository=datatunerx ")
	installChart.Flags().StringVar(&params.ImageSetting.ImagePullPolicy, "image-pull-policy", "", "Specify a image pull policy for the chart version to use. This constraint can be anything that the Helm client supports, for example --image-pull-policy=Always ")
	installChart.Flags().StringVar(&params.ImageSetting.ImagePullSecret, "image-pull-secret", "", "Specify a image pull secret for the chart version to use. This constraint can be anything that the Helm client supports, for example --image-pull-secret=datatunerx ")
	installChart.Flags().BoolVar(&params.ImageSetting.Push, "push", false, "Specify a push for the chart version to use. This constraint can be anything that the Helm client supports, for example --push=true ")
	installChart.Flags().StringVar(&params.ImageSetting.ImageFileDir, "image-file-dir", "", "Specify a image file dir for the chart version to use. This constraint can be anything that the Helm client supports, for example --image-file-dir=/tmp ")
	installChart.Flags().BoolVar(&params.Wait, "wait", false, "Wait for installation to have completed")
	installChart.Flags().DurationVar(&params.WaitDuration, "wait-duration", valueobject.InstallTimeout, "Maximum time to wait for resources to be ready")
	installChart.Flags().BoolVar(&params.DryRun, "dry-run", false, "Simulate an install")
	return installChart
}
