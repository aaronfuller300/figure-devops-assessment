package main

import (
    "context"
    "flag"
    "fmt"
    "os"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "absolute path to the kubeconfig file")
    flag.Parse()

    // Build the config from the kubeconfig file
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        fmt.Printf("Error building kubeconfig: %s\n", err.Error())
        os.Exit(1)
    }

    // Create a clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating Kubernetes client: %s\n", err.Error())
        os.Exit(1)
    }

    // Retrieve Pods
    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing pods: %s\n", err.Error())
        os.Exit(1)
    }

    // Print Pod names
    fmt.Println("Pods in the default namespace:")
    for _, pod := range pods.Items {
        fmt.Println(pod.Name)
    }
}