package datatunerx

import (
	"fmt"

	"github.com/DataTunerX/dtx-ctl/pkg/infrastructure/k8s"
	"github.com/spf13/cobra"
)

var (
	k8sClient   *k8s.Client
	root        *cobra.Command
	contextName string
	namespace   string
)

func NewDataTunerXCommand() *cobra.Command {
	root = &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.HasParent() {
				return nil
			}
			c, err := k8s.NewClient(contextName, "", namespace)
			if err != nil {
				return fmt.Errorf("create k8s client failed, err: %v", err)
			}
			k8sClient = c
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		Use: "dtx-ctl",
	}
	root.PersistentFlags().StringVar(&contextName, "context", "", "Kubernetes configuration context")
	root.PersistentFlags().StringVarP(&namespace, "namespace", "n", "datatunerx-dev", "Namespace datatunerx is running in")
	root.AddCommand(addInstallCommand(), addUnInstallCommand())
	return root
}
