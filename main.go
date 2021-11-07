package main

import (
	"context"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/julienschmidt/httprouter"
)

type ClusterVisHandler struct {
	client *kubernetes.Clientset
}

func main() {

	log.Println("Starting")

	visHandler := &ClusterVisHandler{client: CreateK8sClient()}

	r := httprouter.New()

	// r.HandleFunc("/", visHandler.RenderCluster)
	r.GET("/", visHandler.RenderCluster)
	r.GET("/ns", visHandler.GetNamespaces)
	r.GET("/pods/:ns", visHandler.GetPods)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err.Error())
	}
}

func (ctx *ClusterVisHandler) GetNamespaces(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Fetching namespaces")
	namespaces, err := ctx.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	nss := []string{}

	for _, ns := range namespaces.Items {
		nss = append(nss, ns.ObjectMeta.Name)
	}

	js, err := json.Marshal(nss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (clusterVisHandler *ClusterVisHandler) GetPods(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ns := p.ByName("ns")

	log.Printf("Fetching pods for %s ns\n", ns)
	pods, err := clusterVisHandler.client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	log.Printf("Fetched %s pods for %s ns\n", strconv.Itoa(len(pods.Items)), ns)
	js, err := json.Marshal(pods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ClusterVisHandler *ClusterVisHandler) RenderCluster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err.Error())
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
