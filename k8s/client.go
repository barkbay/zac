package k8s

import (
	"k8s.io/client-go/kubernetes"
)

func NewClientSet() *kubernetes.Clientset {
	// Load K8S cluster configuration
	clusterConfig, err := NewClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Get REST client
	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}
