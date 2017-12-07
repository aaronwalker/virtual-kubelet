package aws

import (
	"log"
	"os"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/virtual-kubelet/virtual-kubelet/manager"
	"github.com/virtual-kubelet/virtual-kubelet/providers"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ECSProvider implements the virtual-kubelet provider interface and communicates with AWS ECS APIs.
type ECSProvider struct {
  ecsClient       *ecs.ECS
	resourceManager *manager.ResourceManager
	nodeName        string
	operatingSystem string
	cluster					string
	region          string
	accessKey       string
	secretKey       string
	cpu             string
	memory          string
	pods            string
}

// NewECSProvider creates a new ECSProvider
func NewECSProvider(config string, rm *manager.ResourceManager, nodeName, operatingSystem string) (*ECSProvider, error) {
	var p ECSProvider
	var err error

	p.resourceManager = rm

	if config != "" {
		f, err := os.Open(config)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if err := p.loadConfig(f); err != nil {
			return nil, err
		}
	}

	if ak := os.Getenv("AWS_ACCESS_KEY_ID"); ak != "" {
		p.accessKey = ak
	}

	if sk := os.Getenv("AWS_SECRET_ACCESS_KEY"); sk != "" {
		p.secretKey = sk
	}

	if r := os.Getenv("AWS_REGION"); r != "" {
		p.region = r
	}

	p.operatingSystem = operatingSystem
	p.nodeName = nodeName

  sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(p.region)},
  )
	p.ecsClient = ecs.New(sess)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// CreatePod accepts a Pod definition and creates
// a ecs service deployment
func (p *ECSProvider) CreatePod(pod *v1.Pod) error {
	log.Println("creating ecs task from pod spec")
	return nil
}

// UpdatePod is a noop, ecs currently does not support live updates of a pod.
func (p *ECSProvider) UpdatePod(pod *v1.Pod) error {
	return nil
}

// DeletePod deletes the specified pod out of ecs.
func (p *ECSProvider) DeletePod(pod *v1.Pod) error {
	log.Println("removing ecs task")
	return nil
}

// GetPod returns a pod by name that is running inside ecs
// returns nil if a pod by that name is not found.
func (p *ECSProvider) GetPod(namespace, name string) (*v1.Pod, error) {
	log.Println("getting ecs task details")
	return nil, nil
}

// GetPodStatus returns the status of a pod by name that is running inside ecs
// returns nil if a pod by that name is not found.
func (p *ECSProvider) GetPodStatus(namespace, name string) (*v1.PodStatus, error) {
	return nil, nil
}

// GetPods returns a list of all pods known to be running within ecs.
func (p *ECSProvider) GetPods() ([]*v1.Pod, error) {
	return nil, nil
}

// Capacity returns a resource list containing the capacity limits set for ecs.
func (p *ECSProvider) Capacity() v1.ResourceList {
	// TODO: These should be configurable
	return v1.ResourceList{
		"cpu":    resource.MustParse("20"),
		"memory": resource.MustParse("100Gi"),
		"pods":   resource.MustParse("20"),
	}
}

// NodeConditions returns a list of conditions (Ready, OutOfDisk, etc), for updates to the node status
// within Kuberentes.
func (p *ECSProvider) NodeConditions() []v1.NodeCondition {
	// TODO: Make these dynamic and augment with custom ecs specific conditions of interest
	return []v1.NodeCondition{
		{
			Type:               "Ready",
			Status:             v1.ConditionTrue,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletReady",
			Message:            "kubelet is ready.",
		},
		{
			Type:               "OutOfDisk",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientDisk",
			Message:            "kubelet has sufficient disk space available",
		},
		{
			Type:               "MemoryPressure",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientMemory",
			Message:            "kubelet has sufficient memory available",
		},
		{
			Type:               "DiskPressure",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasNoDiskPressure",
			Message:            "kubelet has no disk pressure",
		},
		{
			Type:               "NetworkUnavailable",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "RouteCreated",
			Message:            "RouteController created a route",
		},
	}

}

// OperatingSystem returns the operating system for this provider.
// This is a noop to default to Linux for now.
func (p *ECSProvider) OperatingSystem() string {
	return providers.OperatingSystemLinux
}
