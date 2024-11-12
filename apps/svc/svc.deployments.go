package svc

import (
	adapter "deployment-service/apps/repository/adapter"
	"fmt"
	"time"
)

type DeplyomentService struct {
	repository *adapter.Repository
}

func (svc DeplyomentService) GetDeploymentsByNamespace(namespace string) ([]map[string]interface{}, error) {
	deployments, err := svc.repository.Kubernetes.ListDeployments(namespace)
	var deploymentInfo []map[string]interface{}
	if err == nil {

		for _, d := range deployments {
			replicasNeeded := int32(0)
			replicasAvailable := int32(0)
			if d.Spec.Replicas != nil {
				replicasNeeded = *d.Spec.Replicas
			}
			if d.Status.AvailableReplicas > 0 {
				replicasAvailable = d.Status.AvailableReplicas
			}

			// Calculate the age of the deployment
			age := time.Since(d.CreationTimestamp.Time)

			// Fetch resource utilization if you have metrics access (optional, requires Metrics API)
			// For example, get memory and CPU usage per pod (implement separately)
			// resourceUtilization := getDeploymentResourceUtilization(d.Name, namespace, k)

			// Create YAML manifest string
			// deploymentYAML, err := yaml.Marshal(d)
			// if err != nil {
			// 	fmt.Printf("failed to marshal deployment %s to YAML: %v\n", d.Name, err)
			// 	deploymentYAML = []byte("Error generating YAML")
			// }

			deploymentDetails := map[string]interface{}{
				"Deployment Name":    d.Name,
				"Namespace":          d.Namespace,
				"Replicas Needed":    replicasNeeded,
				"Replicas Available": replicasAvailable,
				"Status":             d.Status.Conditions,
				"Age":                age.String(),
				"Creation Time":      d.CreationTimestamp.Time,
				// "Resource Utilization": resourceUtilization, // CPU/Memory data as needed
				// "Deployment Manifest": string(deploymentYAML),
			}

			deploymentInfo = append(deploymentInfo, deploymentDetails)
		}
	}
	// fmt.Println(deploymentInfo)
	// svc.repository.Kubernetes.ListServices("default")
	return deploymentInfo, nil
}

func (svc DeplyomentService) GetTenantKubernetesInfo(namespace string) (map[string]interface{}, error) {
	clusterInfo := make(map[string]interface{})

	// Fetch Namespaces
	namespaces, err := svc.repository.Kubernetes.ListNamespaces()
	if err != nil {
		return nil, err
	}
	var namespaceNames []string
	for _, ns := range namespaces {
		namespaceNames = append(namespaceNames, ns.Name)
	}
	clusterInfo["Namespaces"] = namespaceNames

	// Fetch Pods
	pods, err := svc.repository.Kubernetes.ListPods(namespace)
	if err != nil {
		return nil, err
	}
	var podInfo []map[string]string
	for _, pod := range pods {
		podInfo = append(podInfo, map[string]string{
			"Pod Name":   pod.Name,
			"Namespace":  pod.Namespace,
			"Node Name":  pod.Spec.NodeName,
			"Status":     string(pod.Status.Phase),
			"Start Time": pod.Status.StartTime.Format(time.RFC3339),
		})
	}
	clusterInfo["Pods"] = podInfo

	// Fetch Kubernetes Version
	k8sVersion, err := svc.repository.Kubernetes.GetKubernetesVersion()
	if err != nil {
		return nil, err
	}
	clusterInfo["Kubernetes Version"] = k8sVersion

	// Fetch Deployments
	deployments, err := svc.repository.Kubernetes.ListDeployments(namespace)
	if err != nil {
		return nil, err
	}
	var deploymentInfo []map[string]string
	for _, deploy := range deployments {
		replicas := "0"
		if deploy.Spec.Replicas != nil {
			replicas = fmt.Sprintf("%d", *deploy.Spec.Replicas)
		}
		deploymentInfo = append(deploymentInfo, map[string]string{
			"Deployment Name": deploy.Name,
			"Namespace":       deploy.Namespace,
			"Replicas":        replicas,
			"Available":       fmt.Sprintf("%d", deploy.Status.AvailableReplicas),
			"Status":          deploy.CreationTimestamp.GoString(),
			"Image":           deploy.Spec.Template.Spec.Containers[0].Image,
		})
	}
	clusterInfo["Deployments"] = deploymentInfo

	// Fetch Services
	services, err := svc.repository.Kubernetes.ListServices(namespace)
	if err != nil {
		return nil, err
	}
	var serviceInfo []map[string]string
	for _, svc := range services {
		serviceInfo = append(serviceInfo, map[string]string{
			"Service Name": svc.Name,
			"Namespace":    svc.Namespace,
			"Type":         string(svc.Spec.Type),
			"Cluster IP":   svc.Spec.ClusterIP,
		})
	}
	clusterInfo["Services"] = serviceInfo
	return clusterInfo, nil
}

func (svc DeplyomentService) GetDeploymentByName(namespace, deploymentName string) (map[string]interface{}, error) {
	return svc.repository.Kubernetes.GetDeploymentByName(namespace, deploymentName)
}

func (svc DeplyomentService) CreateNamespace(namespace string) error {
	return svc.repository.Kubernetes.CreateNamespace(namespace)
}
