package ServiceCore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"net/rpc"
	"os"
)

type GalaxyClient struct {
	Version int

	ClientHost string
	ClientPort int

	client *rpc.Client
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
	nameEnv := ""
	portEnv := ""
	if KubernetesUtil.IsRunningInKubernetes() {
		nameEnv = os.Getenv("GALAXY_SERVICE_HOST")
		portEnv = os.Getenv("GALAXY_SERVICE_PORT")
	} else {
		nameEnv = os.Getenv("GALAXY_HOST")
		portEnv = os.Getenv("GALAXY_PORT")

		if nameEnv == "" {
			nameEnv = "localhost"
		}

		if portEnv == "" {
			portEnv = "3000"
		}
	}
	address := fmt.Sprintf("%s:%s", nameEnv, portEnv)

	return NewGalaxyClientWithAddress(address)
}
