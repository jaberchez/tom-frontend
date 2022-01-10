package k8s

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	clientset *kubernetes.Clientset
)

func init() {
	err := createK8sClientSet()

	if err != nil {
		log.Fatal(fmt.Sprintf("unable to create kubernetes clientset: %s", err.Error()))
	}
}

//func GetEndpoints(name string, namespace string) (map[string]int, error) {
//	servers := make(map[string]int)
//
//	// Get Endpoints
//	endpoints, err := clientset.CoreV1().Endpoints(namenamespace).Get(context.Background(),
//		name, metav1.GetOptions{})
//
//	if err != nil {
//		return nil, err
//	}
//
//	if endpoints.Subsets
//
//	return endpoints, nil
//}

func createK8sClientSet() error {
	var config *rest.Config

	// Creates the in-cluster config
	config, err := rest.InClusterConfig()

	if err != nil {
		// Try with kubeconfig out of the cluster
		home := homedir.HomeDir()

		if len(home) == 0 {
			return errors.New("home dir not found")
		}

		kubeconfig := filepath.Join(home, ".kube", "config")

		// Use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

		if err != nil {
			return err
		}
	}

	// Creates the clientset
	clientset, err = kubernetes.NewForConfig(config)

	if err != nil {
		return err
	}

	return nil
}
