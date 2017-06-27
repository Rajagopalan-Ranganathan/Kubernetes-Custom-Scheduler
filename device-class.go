/*
   // Authors:
   Rajagopalan Ranganathan (rajagopalan.ranganthan@aalto.fi)
   Sunil Kumar Mohanty (sunil.mohanty@aalto.fi)

   The following source code, has been created for academic purpose to experiment and use a
   custom Kubernetes container scheduler logic.

   It iterates through the PODs and assigns a Node to POD based on the custom logic

*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
	"time"
)

/*
 PODS Json structure
 It maps the JSON of POD JSON to a strcuture
*/

type Pods struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				KubernetesIoCreatedBy                 string `json:"kubernetes.io/created-by"`
				SchedulerAlphaKubernetesIoCriticalPod string `json:"scheduler.alpha.kubernetes.io/critical-pod"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			GenerateName      string    `json:"generateName"`
			Labels            struct {
				K8SApp                string `json:"k8s-app"`
				PodTemplateGeneration string `json:"pod-template-generation"`
				Network               string `json:"network"`
				Category              string `json:"category"`
			} `json:"labels"`
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
			OwnerReferences []struct {
				APIVersion         string `json:"apiVersion"`
				BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
				Controller         bool   `json:"controller"`
				Kind               string `json:"kind"`
				Name               string `json:"name"`
				UID                string `json:"uid"`
			} `json:"ownerReferences"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			Containers []struct {
				Args    []string `json:"args"`
				Command []string `json:"command"`
				Env     []struct {
					Name      string `json:"name"`
					ValueFrom struct {
						FieldRef struct {
							APIVersion string `json:"apiVersion"`
							FieldPath  string `json:"fieldPath"`
						} `json:"fieldRef"`
					} `json:"valueFrom"`
				} `json:"env"`
				Image           string `json:"image"`
				ImagePullPolicy string `json:"imagePullPolicy"`
				Name            string `json:"name"`
				Resources       struct {
				} `json:"resources"`
				TerminationMessagePath   string `json:"terminationMessagePath"`
				TerminationMessagePolicy string `json:"terminationMessagePolicy"`
				VolumeMounts             []struct {
					MountPath string `json:"mountPath"`
					Name      string `json:"name"`
					ReadOnly  bool   `json:"readOnly,omitempty"`
				} `json:"volumeMounts"`
			} `json:"containers"`
			DNSPolicy    string `json:"dnsPolicy"`
			HostNetwork  bool   `json:"hostNetwork"`
			NodeName     string `json:"nodeName"`
			NodeSelector struct {
				NodeRoleKubernetesIoMaster string `json:"node-role.kubernetes.io/master"`
			} `json:"nodeSelector"`
			RestartPolicy   string `json:"restartPolicy"`
			SchedulerName   string `json:"schedulerName"`
			SecurityContext struct {
			} `json:"securityContext"`
			ServiceAccount                string `json:"serviceAccount"`
			ServiceAccountName            string `json:"serviceAccountName"`
			TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
			Tolerations                   []struct {
				Effect   string `json:"effect,omitempty"`
				Key      string `json:"key"`
				Operator string `json:"operator,omitempty"`
			} `json:"tolerations"`
			Volumes []struct {
				HostPath struct {
					Path string `json:"path"`
				} `json:"hostPath,omitempty"`
				Name   string `json:"name"`
				Secret struct {
					DefaultMode int    `json:"defaultMode"`
					SecretName  string `json:"secretName"`
				} `json:"secret,omitempty"`
			} `json:"volumes"`
		} `json:"spec"`
		Status struct {
			Conditions []struct {
				LastProbeTime      interface{} `json:"lastProbeTime"`
				LastTransitionTime time.Time   `json:"lastTransitionTime"`
				Status             string      `json:"status"`
				Type               string      `json:"type"`
			} `json:"conditions"`
			ContainerStatuses []struct {
				ContainerID string `json:"containerID"`
				Image       string `json:"image"`
				ImageID     string `json:"imageID"`
				LastState   struct {
				} `json:"lastState"`
				Name         string `json:"name"`
				Ready        bool   `json:"ready"`
				RestartCount int    `json:"restartCount"`
				State        struct {
					Running struct {
						StartedAt time.Time `json:"startedAt"`
					} `json:"running"`
				} `json:"state"`
			} `json:"containerStatuses"`
			HostIP    string    `json:"hostIP"`
			Phase     string    `json:"phase"`
			PodIP     string    `json:"podIP"`
			QosClass  string    `json:"qosClass"`
			StartTime time.Time `json:"startTime"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
	} `json:"metadata"`
	ResourceVersion string `json:"resourceVersion"`
	SelfLink        string `json:"selfLink"`
}

/*
 Nodes Json structure
 It maps the JSON of Node JSON to a strcuture
*/

type Nodes struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				NodeAlphaKubernetesIoTTL                         string `json:"node.alpha.kubernetes.io/ttl"`
				VolumesKubernetesIoControllerManagedAttachDetach string `json:"volumes.kubernetes.io/controller-managed-attach-detach"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Labels            struct {
				BetaKubernetesIoArch string `json:"beta.kubernetes.io/arch"`
				BetaKubernetesIoOs   string `json:"beta.kubernetes.io/os"`
				KubernetesIoHostname string `json:"kubernetes.io/hostname"`
				Network              string `json:"network"`
				Category             string `json:"category"`
			} `json:"labels"`
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			ExternalID string `json:"externalID"`
		} `json:"spec"`
		Status struct {
			Addresses []struct {
				Address string `json:"address"`
				Type    string `json:"type"`
			} `json:"addresses"`
			Allocatable struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
				Pods   string `json:"pods"`
			} `json:"allocatable"`
			Capacity struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
				Pods   string `json:"pods"`
			} `json:"capacity"`
			Conditions []struct {
				LastHeartbeatTime  time.Time `json:"lastHeartbeatTime"`
				LastTransitionTime time.Time `json:"lastTransitionTime"`
				Message            string    `json:"message"`
				Reason             string    `json:"reason"`
				Status             string    `json:"status"`
				Type               string    `json:"type"`
			} `json:"conditions"`
			DaemonEndpoints struct {
				KubeletEndpoint struct {
					Port int `json:"Port"`
				} `json:"kubeletEndpoint"`
			} `json:"daemonEndpoints"`
			Images []struct {
				Names     []string `json:"names"`
				SizeBytes int      `json:"sizeBytes"`
			} `json:"images"`
			NodeInfo struct {
				Architecture            string `json:"architecture"`
				BootID                  string `json:"bootID"`
				ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
				KernelVersion           string `json:"kernelVersion"`
				KubeProxyVersion        string `json:"kubeProxyVersion"`
				KubeletVersion          string `json:"kubeletVersion"`
				MachineID               string `json:"machineID"`
				OperatingSystem         string `json:"operatingSystem"`
				OsImage                 string `json:"osImage"`
				SystemUUID              string `json:"systemUUID"`
			} `json:"nodeInfo"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
	} `json:"metadata"`
	ResourceVersion string `json:"resourceVersion"`
	SelfLink        string `json:"selfLink"`
}
type MyJsonName struct {
	Example struct {
		From struct {
			Json bool `json:"json"`
		} `json:"from"`
	} `json:"example"`
}

/*
	Main function
	Runs for ever, sleeps for every one second
*/

func main() {
	for {
		schedulePods()
		time.Sleep(time.Second)
	}

}

/*
	@Function name: postbind
	@Paramin: podname, nodename string
	@return : none, panics during error
	@desc: prepares a JSON string with the Node and POD names and sends the binding request to the master.
	Binds a POD to a NODE
*/

func postbind(podname, nodename string) {

	url := fmt.Sprintf("http://localhost:8001/api/v1/namespaces/default/pods/%s/binding", podname)
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"apiVersion":"v1", "kind": "Binding", "metadata": {"name": "` + podname + `"}, "target": {"apiVersion": "v1", "kind": "Node", "name": "` + nodename + `"}}`)
	fmt.Println(string(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

/*
	@Function name: schedulePods
	@Paramin: <none>
	@return : none, prints error
	@desc: Iterate through the PODs and get POD name, category and network fields and call the assign_node_to_pod function
*/

func schedulePods() {
	app := "kubectl"
	arg0 := "get"
	arg1 := "pods"
	arg2 := "-o"
	arg3 := "json"

	cmd := exec.Command(app, "--server", "localhost:8001", arg0, arg1, arg2, arg3)
	pods_byte, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	var pods_json Pods
	json.Unmarshal(pods_byte, &pods_json)

	for _, elem := range pods_json.Items {
		if elem.Spec.SchedulerName == "my-scheduler" && len(elem.Spec.NodeName) == 0 {
			assign_node_to_pod(elem.Metadata.Name, elem.Metadata.Labels.Category, elem.Metadata.Labels.Network)

		}

	}

}

/*
	@Function name: assign_node_to_pod
	@Paramin: podname, category, network string
	@return : none, prints error
	@desc: Iterate through the Nodes and select a suitable node for the received POD. Check for the labels Category and network and assign
	the Node. When no suitable criterion is found, assign a node randomly

*/

func assign_node_to_pod(podname, category, network string) {

	app := "kubectl"
	arg0 := "get"
	arg1 := "nodes"
	arg2 := "-o"
	arg3 := "json"
	nodename := ""

	cmd := exec.Command(app, "--server", "localhost:8001", arg0, arg1, arg2, arg3)
	nodes_byte, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	var nodes_json Nodes
	json.Unmarshal(nodes_byte, &nodes_json)

	var avl_nodes []string
	for _, elem := range nodes_json.Items {
		if len(network) == 0 {
			elem.Metadata.Labels.Network = ""
		}
		if len(category) == 0 {
			elem.Metadata.Labels.Category = ""
		}
		if network == elem.Metadata.Labels.Network && category == elem.Metadata.Labels.Category {
			avl_nodes = append(avl_nodes, elem.Metadata.Name)
		}
	}
	nodename = avl_nodes[rand.Intn(len(avl_nodes))]
	if len(nodename) != 0 {
		postbind(podname, nodename)
	}

}
