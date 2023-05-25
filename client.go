package GalaxyClient

import "net/rpc"

type GalaxyClient struct {
	Version int
	client  *rpc.Client
}

func NewGalaxyClient(galaxyAddress string) (*GalaxyClient, error) {
	if galaxyAddress == "" {
		galaxyAddress = "127.0.0.1:1234"
	}

	client, err := rpc.DialHTTP("tcp", galaxyAddress)
	if err != nil {
		return nil, err
	}

	return &GalaxyClient{
		Version: 1,
		client:  client,
	}, nil
}
