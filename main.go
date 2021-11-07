package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/julienschmidt/httprouter"
	"github.com/kszpakowski/go-playground/pkg/controller"
)

func main() {

	log.Println("Starting")

	controller := &controller.KubeApiController{Client: CreateK8sClient()}

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "./static/index.html")
	})
	r.GET("/ns", controller.GetNamespaces)
	r.GET("/pods/:ns", controller.GetPods)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err.Error())
	}
}

func CreateK8sClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}
