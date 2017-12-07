package aws

import (
  "crypto/md5"
  "encoding/hex"
  "io/ioutil"
	"log"
  "testing"

  "github.com/aws/aws-sdk-go/service/ecs"
  "github.com/aws/aws-sdk-go/service/ecs/ecsiface"
  "github.com/virtual-kubelet/virtual-kubelet/providers"
  "github.com/virtual-kubelet/virtual-kubelet/manager"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
)

var (
	fakeClient *kubernetes.Clientset
)

func init() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	configOverrides := &clientcmd.ConfigOverrides{}
	// if you want to change override values or bind them to flags, there are methods to help you

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatal("unable to create client config")
	}
	fakeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("unable to create new clientset")
	}
}

// Define a mock struct to be used in your unit tests of myFunc.
type mockECSClient struct {
    ecsiface.ECSAPI
}
func (m *mockECSClient) CreateCluster(input *ecs.CreateClusterInput) (*ecs.CreateClusterOutput, error) {
    // mock response/functionality
    return nil, nil
}

const mockConfig = `
Region = "us-east-1"
Cluster = "mycluster"
CPU = "100"
Memory = "100Gi"
Pods = "20"`

func TestNewECSProvider(t *testing.T) {
    // Setup Test
    rm := manager.NewResourceManager(fakeClient)

    p, err := NewECSProvider(WriteTmpConfig(mockConfig), rm, "fake", providers.OperatingSystemLinux)
    if err != nil || p == nil {
      t.Errorf("failed to create ecs provider %s", err)
    }

}

func WriteTmpConfig(cfg string) string {
  file := "/tmp/" + GetMD5Hash(cfg)
  err := ioutil.WriteFile(file, []byte(cfg), 0644)
  if err != nil {
    log.Fatalf("unable to create tmp config file %s", file)
  }
  return file
}

func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}
