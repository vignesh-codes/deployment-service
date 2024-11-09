package adapter

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NewKubernetes initializes the Kubernetes adapter
func NewKubernetes(client *kubernetes.Clientset) *Kubernetes {
	return &Kubernetes{connection: client}
}

// ListDeployments fetches deployments in the specified namespace
func (k *Kubernetes) ListDeployments(namespace string) ([]appsv1.Deployment, error) {
	deploymentsClient := k.connection.AppsV1().Deployments(namespace)
	deployments, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	return deployments.Items, nil
}

// ListServices fetches services in the specified namespace
func (k *Kubernetes) ListServices(namespace string) ([]corev1.Service, error) {
	servicesClient := k.connection.CoreV1().Services(namespace)
	services, err := servicesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	return services.Items, nil
}

// ListNamespaces fetches all namespaces in the cluster
func (k *Kubernetes) ListNamespaces() ([]corev1.Namespace, error) {
	namespacesClient := k.connection.CoreV1().Namespaces()
	namespaces, err := namespacesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	return namespaces.Items, nil
}

// ListPods fetches all pods in the specified namespace
func (k *Kubernetes) ListPods(namespace string) ([]corev1.Pod, error) {
	podsClient := k.connection.CoreV1().Pods(namespace)
	pods, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return pods.Items, nil
}

// ListNodes fetches all nodes in the cluster
func (k *Kubernetes) ListNodes() ([]corev1.Node, error) {
	nodesClient := k.connection.CoreV1().Nodes()
	nodes, err := nodesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	return nodes.Items, nil
}

// GetKubernetesVersion fetches the Kubernetes server version
func (k *Kubernetes) GetKubernetesVersion() (string, error) {
	versionInfo, err := k.connection.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get Kubernetes version: %w", err)
	}
	return versionInfo.GitVersion, nil
}

// UpdateDeploymentByNameReplicas updates the number of replicas for a specified deployment
func (k *Kubernetes) UpdateDeploymentByNameReplicas(namespace, deploymentName string, replicas int32) error {
	deploymentsClient := k.connection.AppsV1().Deployments(namespace)
	deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	deployment.Spec.Replicas = &replicas
	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update replicas for deployment %s: %w", deploymentName, err)
	}

	fmt.Printf("Successfully updated replicas for deployment %s to %d\n", deploymentName, replicas)
	return nil
}

// UpdateDeploymentByNameImageVersion updates the image version for a specified deployment
func (k *Kubernetes) UpdateDeploymentByNameImageVersion(namespace, deploymentName, containerName, newImageVersion string) error {
	deploymentsClient := k.connection.AppsV1().Deployments(namespace)
	deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Update the image for the specified container in the deployment
	updated := false
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			deployment.Spec.Template.Spec.Containers[i].Image = newImageVersion
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("container %s not found in deployment %s", containerName, deploymentName)
	}

	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update image version for deployment %s: %w", deploymentName, err)
	}

	fmt.Printf("Successfully updated image for container %s in deployment %s to %s\n", containerName, deploymentName, newImageVersion)
	return nil
}

// DeleteDeploymentByName deletes a specified deployment from a namespace
func (k *Kubernetes) DeleteDeploymentByName(namespace, deploymentName string) error {
	deploymentsClient := k.connection.AppsV1().Deployments(namespace)
	err := deploymentsClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment %s: %w", deploymentName, err)
	}

	fmt.Printf("Successfully deleted deployment %s from namespace %s\n", deploymentName, namespace)
	return nil
}
