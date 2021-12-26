package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	clientset *kubernetes.Clientset
)

type Backend struct {
	Name string
	IP   string
	Port int
}

func init() {
	err := createK8sClientSet()

	if err != nil {
		log.Fatal(fmt.Sprintf("unable to create kubernetes clientset: %s", err.Error()))
	}
}

func Get(backendService string, namespace string) (backends []Backend, err error) {
	// Get Kubernetes Endpoints belong to this Service
	endpoints, err := clientset.CoreV1().Endpoints(namespace).Get(context.Background(),
		backendService, metav1.GetOptions{})

	if err != nil {
		return
	}

	if err != nil {
		return
	}

	if len(endpoints.Subsets) == 0 {
		err = errors.New("found empty endpoints")
		return
	}

	for _, item := range endpoints.Subsets {
		for _, a := range item.Addresses {
			backend := Backend{
				Name: a.TargetRef.Name,
				IP:   a.IP,
				Port: int(item.Ports[0].Port),
			}

			backends = append(backends, backend)
		}
	}

	return
}

func GetEnvVars(ip string, port int) (envs map[string]string, err error) {
	envs = make(map[string]string)

	timeout := time.Duration(5 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/env", ip, port), nil)

	if err != nil {
		return
	}

	// Appending to existing query args
	//q := req.URL.Query()
	//q.Add("foo", "bar")

	// assign encoded query string to http request
	//req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		//return envs, fmt.Errorf("errored when sending request to the server %s: %s", server, err.Error())
		return
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(responseBody, &envs)

	if err != nil {
		return
	}

	//fmt.Println(resp.Status)
	//fmt.Println(string(responseBody))

	return
}

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
