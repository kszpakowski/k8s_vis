package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeApiController struct {
	Client *kubernetes.Clientset
}

func (controller *KubeApiController) GetNamespaces(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Fetching namespaces")
	namespaces, err := controller.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Fetched %s namespaces\n", strconv.Itoa(len(namespaces.Items)))

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

func (controller *KubeApiController) GetPods(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ns := p.ByName("ns")

	log.Printf("Fetching pods for %s ns\n", ns)
	pods, err := controller.Client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	log.Printf("Fetched %s pods for %s ns\n", strconv.Itoa(len(pods.Items)), ns)
	js, err := json.Marshal(pods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
