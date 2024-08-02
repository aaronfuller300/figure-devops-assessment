package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/retry"
)

func main() {
    kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "absolute path to the kubeconfig file")
    flag.Parse()

    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        fmt.Printf("Error building kubeconfig: %s\n", err.Error())
        os.Exit(1)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating Kubernetes client: %s\n", err.Error())
        os.Exit(1)
    }

    deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing deployments: %s\n", err.Error())
        os.Exit(1)
    }

    for _, deployment := range deployments.Items {
        if containsDatabase(deployment.Name) {
            fmt.Printf("Redeploying deployment: %s\n", deployment.Name)

            err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
                deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
                _, updateErr := clientset.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
                return updateErr
            })

            if err != nil {
                fmt.Printf("Error redeploying deployment %s: %s\n", deployment.Name, err.Error())
            }
        }
    }
}

func containsDatabase(name string) bool {
    return len(name) > 0 && (len(name) >= 8 && name[len(name)-8:] == "database")
}