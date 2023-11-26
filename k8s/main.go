package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
)

func main() {

	//configPath := "/root/.kube/config"
	configPath := "D:/Dev/go/config"
	clientset := createClient(configPath)

	// listPods(clientset)
	// fmt.Printf("pods: %v\n", pods)

	// podName := "myweb-7b9b5bc894-2qr4j"
	// delPod(clientset, podName)

}

func createClient(configPath string) *kubernetes.Clientset {
	var err error
	var kubeconfig []byte
	var clientconfig *rest.Config
	var clientset *kubernetes.Clientset
	if kubeconfig, err = ioutil.ReadFile(configPath); err != nil {
		panic(err)
	}

	if clientconfig, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		panic(err)
	}

	if clientset, err = kubernetes.NewForConfig(clientconfig); err != nil {
		panic(err)
	}
	return clientset
}

func createDynamicClient(configPath string) *dynamic.DynamicClient {
	var err error
	var inConfig *rest.Config
	var dd *dynamic.DynamicClient
	if inConfig, err = clientcmd.BuildConfigFromFlags("", configPath); err != nil {
		panic(err)
	}

	if dd, err = dynamic.NewForConfig(inConfig); err != nil {
		panic(err)
	}
	return dd
}

func createDeploy(yamlFile string) *apps_v1.Deployment {
	var err error
	buf, _ := ioutil.ReadFile(yamlFile)
	var deployment = &apps_v1.Deployment{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode(buf, nil, deployment)
	if err != nil {
		log.Fatalf("Decode err %v", err)
	}
	log.Printf("deploy is %s", deployment.Name)
	return deployment
}

func createDeployment(clientset *kubernetes.Clientset, dd *dynamic.DynamicClient, yamlFile string, namespace string) {
	var err error
	var filebytes []byte
	if filebytes, err = ioutil.ReadFile(yamlFile); err != nil {
		panic(err)
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatal(err)
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			log.Fatal(err)
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatal(err)
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(namespace)
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		obj2, err := dri.Create(context.Background(), unstructuredObj, meta_v1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s/%s created", obj2.GetKind(), obj2.GetName())
	}

}

func listPods(clientset *kubernetes.Clientset) *core_v1.PodList {
	var err error
	var pods *core_v1.PodList
	// 查询pod列表
	if pods, err = clientset.CoreV1().Pods("default").List(context.TODO(), meta_v1.ListOptions{}); err != nil {
		panic(err)
	}
	return pods
}

func delPod(clientset *kubernetes.Clientset, podName string) {
	var err error
	// 删除pod
	if err = clientset.CoreV1().Pods("default").Delete(context.TODO(), podName, meta_v1.DeleteOptions{}); err != nil {
		panic(err)
	}
}

func delDeployment(clientset *kubernetes.Clientset, deploymentName string) {
	var err error
	// 删除deployment
	if err = clientset.AppsV1().Deployments("default").Delete(context.TODO(), deploymentName, meta_v1.DeleteOptions{}); err != nil {
		panic(err)
	}
}
