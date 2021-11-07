package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Namespace struct {
	Id   string
	Name string
}

type Pod struct {
	Id        string
	Name      string
	Namespace string
}

type ClusterData struct {
	Namespaces []Namespace
	Pods       []Pod
}

type ClusterVisHandler struct {
	client *kubernetes.Clientset
}

func main() {

	log.Println("Starting")

	visHandler := &ClusterVisHandler{client: CreateK8sClient()}
	http.HandleFunc("/", visHandler.RenderCluster)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err.Error())
	}
}

func (ClusterVisHandler *ClusterVisHandler) RenderCluster(w http.ResponseWriter, r *http.Request) {
	cluster := ClusterData{}

	namespaces, err := ClusterVisHandler.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, ns := range namespaces.Items {
		cluster.Namespaces = append(cluster.Namespaces, Namespace{"ns_" + ns.Name, ns.Name})
	}

	pods, err := ClusterVisHandler.client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	for _, pod := range pods.Items {
		cluster.Pods = append(cluster.Pods, Pod{"pod_" + pod.Name, pod.Name, pod.Namespace})
	}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, cluster); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
