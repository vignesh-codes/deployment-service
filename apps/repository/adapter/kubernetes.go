package adapter

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) ListDeployments(namespace string) error {
	deploymentsClient := k.connection.AppsV1().Deployments(namespace)
	deployments, err := deploymentsClient.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	for _, d := range deployments.Items {
		fmt.Println("Deployment Name: %s, Replicas: %d\n", d.Name, *d.Spec.Replicas)
	}
	return nil
}
