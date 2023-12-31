package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

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
)

func main() {

	//configPath := "/root/.kube/config"
	configPath := "D:\\Dev\\go\\gostudy\\client-go_samples\\config"
	clientset := createClient(configPath)
	dd := createDynamicClient(configPath)
	// listPods(clientset)
	// fmt.Printf("pods: %v\n", pods)

	// podName := "myweb-7b9b5bc894-2qr4j"
	// delPod(clientset, podName)

	deleteDeployment(clientset, dd, "D:\\Dev\\go\\gostudy\\client-go_samples\\020_deploy_nginx.yaml", "default")
	createDeployment(clientset, dd, "D:\\Dev\\go\\gostudy\\client-go_samples\\020_deploy_nginx.yaml", "default", false, 2)

	// nodes := getNodes(clientset)
	// for _, node := range nodes.Items {
	// 	if node.GetName() == "k8s-node2" {
	// 		fmt.Printf("NodeName: %s\n", node.GetName())
	// 		addTaint(clientset, node, "err", "gpu-disabled")
	// 	}
	// }

	// podList := listPods(clientset, "default")
	// bytes, _ := json.Marshal(podList)
	// fmt.Println(string(bytes))
	// for _, p := range podList.Items {
	// 	fmt.Printf("%v\n", p.GetName())
	// 	delPod(clientset, p.GetName(), "default")
	// }

	// nodes := getNodes(clientset)
	// for _, node := range nodes.Items {
	// 	if node.GetName() == "k8s-node2" {
	// 		fmt.Printf("NodeName: %s\n", node.GetName())
	// 		delTaint(clientset, node, "err", "gpu-disabled")
	// 	}
	// }

	// node := getNode(clientset, "k8s-node2")

	// fmt.Printf("node: %v\n", node)
}

func delTaint(clientset *kubernetes.Clientset, node core_v1.Node, key string, val string) {
	for i, taint := range node.Spec.Taints {
		if taint.Key == key && taint.Value == val {
			node.Spec.Taints = append(node.Spec.Taints[:i], node.Spec.Taints[i+1:]...)
			break
		}
	}
	_, err := clientset.CoreV1().Nodes().Update(context.TODO(), &node, meta_v1.UpdateOptions{})
	if err != nil {
		log.Fatalf("update node failed: %v", err)
	}
	log.Println("Taint deleted successfully")
}

func hasTaint(taints []core_v1.Taint, key string, val string) bool {
	hasTaint := false
	for _, taint := range taints {
		if taint.Key == key && taint.Value == val {
			hasTaint = true
			break
		}
	}
	return hasTaint
}

func addTaint(clientset *kubernetes.Clientset, node core_v1.Node, key string, val string) {
	if !hasTaint(node.Spec.Taints, key, val) {
		taint := &core_v1.Taint{
			Key:    key,
			Value:  val,
			Effect: core_v1.TaintEffectNoSchedule,
		}
		node.Spec.Taints = append(node.Spec.Taints, *taint)

		_, err := clientset.CoreV1().Nodes().Update(context.TODO(), &node, meta_v1.UpdateOptions{})

		if err != nil {
			log.Fatalf("update node failed: %v", err)
		}
		log.Printf("add taint{%v:%v} successful\n", key, val)
	} else {
		log.Println("taint existed")
	}
}

func getNode(clientset *kubernetes.Clientset, nodeName string) *core_v1.Node {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, meta_v1.GetOptions{})
	if err != nil {
		log.Fatalf("get node failed: %v", err)
	}
	return node
}

func getNodes(clientset *kubernetes.Clientset) *core_v1.NodeList {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("get node list failed: %v", err)
	}
	return nodes
}

func createClient(configPath string) *kubernetes.Clientset {
	var err error
	var kubeconfig []byte
	var clientconfig *rest.Config
	var clientset *kubernetes.Clientset
	if kubeconfig, err = ioutil.ReadFile(configPath); err != nil {
		log.Fatalf("read kubeconfig file failed: %v", err)
	}

	if clientconfig, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		log.Fatalf("rest kubeconfig failed: %v", err)
	}

	if clientset, err = kubernetes.NewForConfig(clientconfig); err != nil {
		log.Fatalf("create clientset failed: %v", err)
	}
	return clientset
}

func createDynamicClient(configPath string) *dynamic.DynamicClient {
	var err error
	var inConfig *rest.Config
	var dynamicClient *dynamic.DynamicClient
	if inConfig, err = clientcmd.BuildConfigFromFlags("", configPath); err != nil {
		log.Fatalf("read config file failed: %v", err)
	}

	if dynamicClient, err = dynamic.NewForConfig(inConfig); err != nil {
		log.Fatalf("create dynamic client failed: %v", err)
	}
	return dynamicClient
}

// func createDeploy(yamlFile string) *apps_v1.Deployment {
// 	var err error
// 	buf, _ := ioutil.ReadFile(yamlFile)
// 	var deployment = &apps_v1.Deployment{}
// 	_, _, err = scheme.Codecs.UniversalDeserializer().Decode(buf, nil, deployment)
// 	if err != nil {
// 		log.Fatalf("Decode err %v", err)
// 	}
// 	log.Printf("deploy is %s", deployment.Name)
// 	return deployment
// }

func createDeployment(clientset *kubernetes.Clientset, dynamicClient *dynamic.DynamicClient, yamlFile string, namespace string, needScale bool, newReplicas float64) {
	var err error
	var filebytes []byte
	if filebytes, err = ioutil.ReadFile(yamlFile); err != nil {
		log.Fatalf("read yaml file failed: %v", err)
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}
		if needScale {
			rawObj = changeReplacas(rawObj, newReplicas)
		}
		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatalf("convert to Unstructured failed: %v", err)
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			log.Fatalf("get group resources failed: %v", err)
		}
		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatalf("mapping failed: %v", err)
		}
		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(namespace)
			}
			dri = dynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dynamicClient.Resource(mapping.Resource)
		}

		obj2, err := dri.Create(context.Background(), unstructuredObj, meta_v1.CreateOptions{})
		if err != nil {
			log.Fatalf("create dynamic resource failed: %v", err)
		}
		log.Printf("%s/%s created\n", obj2.GetKind(), obj2.GetName())
	}
}

func changeReplacas(rawObj runtime.RawExtension, newReplicas float64) runtime.RawExtension {
	var b []byte
	b, err := rawObj.MarshalJSON()
	if err != nil {
		log.Fatalf("marshal raw to json faided: %v", err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalf("unmarshal to json faided: %v", err)
	}
	if spec, ok := data["spec"].(map[string]interface{}); ok {
		if replicas, ok := spec["replicas"].(float64); ok {
			log.Printf("scale replicas from %v to %v\n", replicas, newReplicas)
			spec["replicas"] = newReplicas
		}
	}
	b, err = json.Marshal(data)
	if err != nil {
		log.Fatalf("json marshal failed: %v", err)
	}
	rawObj.UnmarshalJSON(b)
	return rawObj
}

func deleteDeployment(clientset *kubernetes.Clientset, dynamicClient *dynamic.DynamicClient, yamlFile string, namespace string) {
	var err error
	var filebytes []byte
	if filebytes, err = ioutil.ReadFile(yamlFile); err != nil {
		log.Fatalf("read yaml file failed: %v", err)
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
			log.Fatalf("convert to Unstructured failed: %v", err)
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			log.Fatalf("get group resources failed: %v", err)
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatalf("mapping failed: %v", err)
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(namespace)
			}
			dri = dynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dynamicClient.Resource(mapping.Resource)
		}

		err = dri.Delete(context.Background(), unstructuredObj.GetName(), meta_v1.DeleteOptions{})
		if err != nil {
			log.Fatalf("delete dployment failed: %v", err)
		}
		log.Printf("%v/%s deleted\n", unstructuredObj.GetKind(), unstructuredObj.GetName())
	}
}

func listPods(clientset *kubernetes.Clientset, namespace string) *core_v1.PodList {
	var err error
	var pods *core_v1.PodList
	// 查询pod列表
	if pods, err = clientset.CoreV1().Pods(namespace).List(context.TODO(), meta_v1.ListOptions{}); err != nil {
		log.Fatalf("get pod list failed: %v", err)
	}
	return pods
}

func delPod(clientset *kubernetes.Clientset, podName string, namespace string) {
	var err error
	// 删除pod
	if err = clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, meta_v1.DeleteOptions{}); err != nil {
		log.Fatalf("delete pod failed: %v", err)
	}
}

func delDeployment(clientset *kubernetes.Clientset, deploymentName string) {
	var err error
	// 删除deployment
	if err = clientset.AppsV1().Deployments("default").Delete(context.TODO(), deploymentName, meta_v1.DeleteOptions{}); err != nil {
		log.Fatalf("delete deployment failed: %v", err)
	}
}
