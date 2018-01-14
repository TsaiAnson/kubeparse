package main

import (
    "fmt"
    "strconv"
    "path/filepath"

    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/util/retry"
)

// Returns an out-of-cluster client using client-go
func getClientSetOut() *kubernetes.Clientset {
    // Getting deploymentsClient object
    home := homedir.HomeDir();
    abspath := filepath.Join(home, ".kube", "config")
    config, err := clientcmd.BuildConfigFromFlags("", abspath)
    if err != nil {
        panic(err)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    return clientset
}

// Returns an in-cluster client using client-go
func getClientSetIn() *kubernetes.Clientset {
    config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

    return clientset
}

// Updates the replica count of Deployment METANAME by the value QUANTITY
func replicaUpdate(clientset *kubernetes.Clientset, metaname string, quantity string) {
    // Getting deployments
    deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

    // Updating deployment
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        // Retrieve the latest version of Deployment before attempting update
        // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
        result, getErr := deploymentsClient.Get(metaname, metav1.GetOptions{})
        if getErr != nil {
            panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
        }

        fmt.Printf("Updating replica count of %v by %v\n", metaname, quantity)

        // Parsing quantity to int32
        i, err := strconv.ParseInt(quantity, 10, 32)
        if err != nil {
            panic(err)
        }

        // Modify replica count
        oldRep := result.Spec.Replicas
        result.Spec.Replicas = int32Ptr(*oldRep + int32(i))
        if *result.Spec.Replicas < int32(1) {
            result.Spec.Replicas = int32Ptr(1)
        }
        _, updateErr := deploymentsClient.Update(result)
        return updateErr
    })
    if retryErr != nil {
        panic(fmt.Errorf("Update failed: %v", retryErr))
    }
    fmt.Printf("Updated replica count of Deployment %v\n", metaname)
}

// Adds a label to node NODE
func addLabel(clientset *kubernetes.Clientset, nodename string, labelkey string, labelvalue string) {
    // Getting node
    nodesClient := clientset.Core().Nodes()

    // Updating node
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        result, getErr := nodesClient.Get(nodename, metav1.GetOptions{})
        if getErr != nil {
            panic(fmt.Errorf("Failed to get latest version of Node: %v", getErr))
        }

        fmt.Printf("Adding label %v:%v to node %v\n", labelkey, labelvalue, nodename)

        // Modify labels
        result.Labels[labelkey] = labelvalue
        _, updateErr := nodesClient.Update(result)
        return updateErr
    })
    if retryErr != nil {
        panic(fmt.Errorf("Update failed: %v", retryErr))
    }
    fmt.Printf("Updated labels of node %v\n", nodename)
}

// Adds a nodeSelector to pod POD* (INCOMPLETE)
func addNodeSel(clientset *kubernetes.Clientset, podname string, labelkey string, labelvalue string) {
    // Getting pod
    podsClient := clientset.Core().Pods(apiv1.NamespaceDefault)


    // Since the nodeSelector field of a POD cannot be modified, we will have to
    // delete and recreate the existing POD with the new nodeSelector
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        result, getErr := podsClient.Get(podname, metav1.GetOptions{})
        if getErr != nil {
            panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
        }

        fmt.Printf("Adding nodeSelector %v:%v to pod %v\n", labelkey, labelvalue, podname)

        result.Spec.NodeName = labelkey
        _, updateErr := podsClient.Update(result)
        return updateErr

        // // Modify labels
        // if result.Spec.NodeSelector == nil {
        //     result.Spec.NodeSelector = make(map[string]string)
        // }
        // result.Spec.NodeSelector[labelkey] = labelvalue
        //
        // deleteErr := podsClient.Delete(podname, &metav1.DeleteOptions{})
        // if deleteErr != nil {
        //     panic(fmt.Errorf("Failed to delete oldPod: %v", deleteErr))
        // }
        //
        // _, createErr := podsClient.Create(result)
        // return createErr
    })
    if retryErr != nil {
        panic(fmt.Errorf("Update failed: %v", retryErr))
    }
    fmt.Printf("Updated nodeSelectors of pod %v\n", podname)
}

func int32Ptr(i int32) *int32 { return &i }

func main() {
    // Unit test for replicaUpdate
    testclient := getClientSetOut()
    //replicaUpdate(testclient, "frontend", "-1")
    //addLabel(testclient, "ip-172-20-38-51.us-west-1.compute.internal", "test", "8")
    addNodeSel(testclient, "frontend-1768566195-ct0qb", "ip-172-20-46-58.us-west-1.compute.internal", "8")
}
