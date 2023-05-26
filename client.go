package ServiceCore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"net/rpc"
	"os"
)

type GalaxyClient struct {
	Version int
	client  *rpc.Client
}

func NewGalaxyClientWithAddress(galaxyAddress string) (*GalaxyClient, error) {
	println("Connecting to Galaxy at ", galaxyAddress)
	client, err := rpc.DialHTTP("tcp", galaxyAddress)
	if err != nil {
		return nil, err
	}

	return &GalaxyClient{
		Version: 1,
		client:  client,
	}, nil
}

func NewGalaxyClient() (*GalaxyClient, error) {
	address := ""
	if KubernetesUtil.IsRunningInKubernetes() {
		nameEnv := os.Getenv("GALAXY_SERVICE_HOST")
		portEnv := os.Getenv("GALAXY_SERVICE_PORT")
		address = fmt.Sprintf("%s:%s", nameEnv, portEnv)
	}

	return NewGalaxyClientWithAddress(address)
}
