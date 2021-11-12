package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	_ "embed"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/julienschmidt/httprouter"
	"github.com/kszpakowski/go-playground/pkg/controller"
)

//go:embed html/index.html
var html []byte

func main() {

	outCluster := flag.Bool("out-cluster", false, "run out cluster configuration, in cluster is default")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()
	if *outCluster {
		log.Printf("Runnin out-cluster mode")
	} else {
		log.Printf("Runnin in-cluster mode")
	}

	k8sApiClient := CreateK8sClient(outCluster, kubeconfig)
	controller := &controller.KubeApiController{Client: k8sApiClient}

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write(html)
	})
	r.GET("/ns", controller.GetNamespaces)
	r.GET("/pods/:ns", controller.GetPods)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err.Error())
	}
}

func CreateK8sClient(outCluster *bool, kubeconfig *string) *kubernetes.Clientset {
	var config rest.Config

	if *outCluster {
		config = *CreateOutClusterK8sApiConfig(kubeconfig)
	} else {
		config = *CreateInClusterK8sApiConfig()
	}

	// create the clientset
	client, err := kubernetes.NewForConfig(&config)
	if err != nil {
		panic(err.Error())
	}

	return client
}

func CreateInClusterK8sApiConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

func CreateOutClusterK8sApiConfig(kubeconfig *string) *rest.Config {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}
