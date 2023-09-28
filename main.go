package main

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	//"k8s.io/client-go/util/kubeconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	//"k8s.io/client-go/util/clientcmd"
	"k8s.io/client-go/tools/clientcmd"
)

// Custom Resource
var ck3sResource = schema.GroupVersionResource{
	Group:   "stack.civo.com",
	Version: "v1alpha1",
	//Resource: "CivoK3sCluster",
	Resource: "civok3sclusters",
}

func main() {

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// CONNECT TO KUBERNETES ********************************************************
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// try to fetch the config from the cluster if this app is running into a pod with the appropriate ServiceAccount
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
		// try to fetch the config from the ~/.kube/config if it's running outside the cluster
		//kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config_kubefixtures")
		config, err = clientcmd.BuildConfigFromFlags("", "/Users/alessandroargentieri/.kube/config_kubefixtures")
		if err != nil {
			panic(err.Error())
		}
	}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// FETCH THE STANDARD RESOURCES (pods, deployments, services, ...)WITH A STANDARD CLIENT (kubernetes.ClientSet)
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	// kubernetes standard client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// get all Pods in all namespaces:
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//panic(err.Error())
		fmt.Println(err)
	} else {
		for _, pod := range pods.Items {
			fmt.Println(pod)
		}
	}

	fmt.Println()
	fmt.Println()

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// FETCH THE CUSTOM RESOURCES (civok3sclusters, ...) WITHOUT A GENERATED CLIENT, USING THE DYNAMIC CLIENT (dynamic.DynamicClient)
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	// kubernetes dynamic client: to interact with Custom Resource without having specific client SDK for the Custom Resource:
	// the client is normally generated through:
	//    $ go get k8s.io/code-generator/cmd/client-gen
	//    $ client-gen --clientset-name versioned --input your-api-group.example.com/v1
	dynamicClient := dynamic.NewForConfigOrDie(config)

	result, err := dynamicClient.Resource(ck3sResource).Namespace("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ck3sUnstructured := range result.Items {

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// print the raw (unstructured.Unstructured) data *************************
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		fmt.Printf("Fetched Civo K3s Cluster (of raw type unstructured.Unstructured): %+v\n", ck3sUnstructured)

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// access some specific spec from the raw data ****************************
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		k3sVersion, found, err := unstructured.NestedString(ck3sUnstructured.UnstructuredContent(), "spec", "version")
		if err != nil {
			fmt.Println(err)
		}
		if !found {
			fmt.Println("k3s version field not found on CivoK3sCluster spec")
		}
		fmt.Printf("k3s version: %s", k3sVersion)

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// patch the spec via the unstructured.Unstructured raw interface *********
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		patch := []interface{}{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/spec/masters",
				"value": 2,
			},
		}

		payload, err := json.Marshal(patch)
		if err != nil {
			fmt.Println(err)
		}

		_, err = dynamicClient.Resource(ck3sResource).Namespace(ck3sUnstructured.GetNamespace()).Patch(context.TODO(), ck3sUnstructured.GetName(), types.JSONPatchType, payload, metav1.PatchOptions{})
		if err != nil {
			fmt.Println(err)
		}

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// convert the raw data into a specific struct ****************************
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		ck3s := CivoK3sCluster{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(ck3sUnstructured.UnstructuredContent(), &ck3s)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Fetched Civo K3s Cluster (of type CivoK3sCluster): %+v\n", ck3s)

		ck3s.Spec.MasterSize = &[]string{"veery-big"}[0]

	}

	/*crName := "test-pools"
	crNamespace := "cust-default-9abdbb75-7505e89134a2"

	result, err := customResource.Namespace(crNamespace).Get(context.TODO(), crName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Fetched Custom Resource: %+v\n", result)
	*/
}
