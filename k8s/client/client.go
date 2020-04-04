package client

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sync"
)

func loadFromKubeConfig() (*rest.Config, error) {
	fmt.Println("Attempt to load from config")
	var kubeConfigPath *string
	kubeConfigPath = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "kube config file")
	flag.Parse()
	restCfg, restErr := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	return restCfg, restErr
}

func assumeServiceAccountAccess() (*rest.Config, error) {
	fmt.Println("attempt to load from serviceaccount")
	return rest.InClusterConfig()
}

func insideKube() bool {
	apiServerHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	apiServerPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	return len(apiServerHost) > 0 && len(apiServerPort) > 0
}

var clientSet *kubernetes.Clientset
var mut = sync.Mutex{}

func GetClient() *kubernetes.Clientset {
	if clientSet != nil {
		return clientSet
	}
	initialize()
	return clientSet
}

func initialize() {
	mut.Lock()
	var cfg *rest.Config
	var err error
	if insideKube() {
		cfg, err = assumeServiceAccountAccess()
	} else {
		cfg, err = loadFromKubeConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	cset, cerr := kubernetes.NewForConfig(cfg)
	if cerr != nil {
		panic(cerr.Error())
	}
	clientSet = cset
	mut.Unlock()
}
