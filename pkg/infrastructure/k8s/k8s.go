package k8s

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Client struct {
	Clientset        kubernetes.Interface
	RawConfig        clientcmdapi.Config
	HelmActionConfig *action.Configuration
	context          string
}

func NewClient(contextName, kubeConfig, namespace string) (*Client, error) {
	restClientGetter := genericclioptions.ConfigFlags{
		Context:    &contextName,
		KubeConfig: &kubeConfig,
	}
	rawKubeConfigLoader := restClientGetter.ToRawKubeConfigLoader()
	config, err := rawKubeConfigLoader.ClientConfig()
	if err != nil {
		return nil, err
	}
	rawConfig, err := rawKubeConfigLoader.RawConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	if contextName == "" {
		contextName = rawConfig.CurrentContext
	}

	helmDriver := ""
	actionConfig := action.Configuration{}
	logger := func(format string, v ...interface{}) {}
	if err := actionConfig.Init(&restClientGetter, namespace, helmDriver, logger); err != nil {
		return nil, err
	}

	return &Client{
		Clientset:        clientset,
		RawConfig:        rawConfig,
		context:          contextName,
		HelmActionConfig: &actionConfig,
	}, nil
}

// todo(tigerK)
func (c *Client) GetDataTunerXStatus() {
}

func (c *Client) DeleteDeployment(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}

func (c *Client) DeleteService(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}

func (c *Client) ListDeploymentByLabel(ctx context.Context, namespace string, label map[string]string) (*v1.DeploymentList, error) {
	return c.Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteDeploymentByLabel(ctx context.Context, namespace string, label map[string]string) error {
	return c.Clientset.AppsV1().Deployments(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) ListServiceByLabel(ctx context.Context, namespace string, label map[string]string) (*corev1.ServiceList, error) {
	return c.Clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteServiceAccountByLabel(ctx context.Context, namespace string, label map[string]string) error {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) ListClusterRoleByLabel(ctx context.Context, label map[string]string) (*rbacv1.ClusterRoleList, error) {
	return c.Clientset.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteClusterRoleByLabel(ctx context.Context, label map[string]string) error {
	return c.Clientset.RbacV1().ClusterRoles().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteCLusterRoleBindingByLabel(ctx context.Context, label map[string]string) error {
	return c.Clientset.RbacV1().ClusterRoleBindings().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteSecretByLabel(ctx context.Context, namespace string, label map[string]string) error {
	return c.Clientset.CoreV1().Secrets(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteMutatingWebhookConfigurationsByLabel(ctx context.Context, label map[string]string) error {
	return c.Clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}

func (c *Client) DeleteValidatingWebhookConfigurationsByLabel(ctx context.Context, label map[string]string) error {
	return c.Clientset.AdmissionregistrationV1().ValidatingWebhookConfigurations().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labels.Set(label).String()})
}
