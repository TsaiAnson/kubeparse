package main

import (
    "fmt"

    appsv1beta1 "k8s.io/api/apps/v1beta1"
    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/retry"

    //"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func replicaUpdate(abspath string, metaname string, magnitude string) {
    // Getting deploymentsClient object
    config, err := clientcmd.BuildConfigFromFlags("", abspath)
    if err != nil {
        panic(err)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

    // Updating deployment
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        // Retrieve the latest version of Deployment before attempting update
        // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
        result, getErr := deploymentsClient.Get(metaname, metav1.GetOptions{})
        if getErr != nil {
            panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
        }

        fmt.Println("Updating replica count of %v by %v", metaname, magnitude)

        // Parsing magnitude to int32
        i, err := strconv.ParseInt(magnitude, 10, 32)
        if err != nil {
            panic(err)
        }

        // Modify replica count
        oldRep := result.Spec.Replicas
        result.Spec.Replicas = oldRep + int32(i)
        if result.Spec.Replicas < 1 {
            result.Spec.Replicas = 1
        }
        _, updateErr := deploymentsClient.Update(result)
        return updateErr
    })
    if retryErr != nil {
        panic(fmt.Errorf("Update failed: %v", retryErr))
    }
    fmt.Println("Updated replica count of Deployment %v from %v to %v", metaname, oldRep, result.Spec.Replicas)
}

func replicaUpdatev1(metaname string, magnitude string) {
    config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"
	deploymentsClient := clientset.Extensions().Deployments(namespace)

    // Updating deployment
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        // Retrieve the latest version of Deployment before attempting update
        // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
        result, getErr := deploymentsClient.Get(metaname, metav1.GetOptions{})
        if getErr != nil {
            panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
        }

        fmt.Println("Updating replica count of %v by %v", metaname, magnitude)

        // Parsing magnitude to int32
        i, err := strconv.ParseInt(magnitude, 10, 32)
        if err != nil {
            panic(err)
        }

        // Modify replica count
        oldRep := result.Spec.Replicas
        result.Spec.Replicas = oldRep + int32(i)
        if result.Spec.Replicas < 1 {
            result.Spec.Replicas = 1
        }
        _, updateErr := deploymentsClient.Update(result)
        return updateErr
    })
    if retryErr != nil {
        panic(fmt.Errorf("Update failed: %v", retryErr))
    }
    fmt.Println("Updated replica count of Deployment %v from %v to %v", metaname, oldRep, result.Spec.Replicas)
}
