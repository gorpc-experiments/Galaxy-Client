package ServiceCore

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/rpc"
)

func (galaxy *GalaxyClient) Call(method string, args any, reply any) error {
	serviceUrl, err := galaxy.LookUp(method)

	if err != nil {
		log.Println("Unable to find a service for method: ", method, err)
	}

	client, err := rpc.DialHTTP("tcp", serviceUrl)

	if err != nil {
		log.Println("Unable to call service for method:", method, err)
	}

	err = client.Call(method, args, reply)

	spew.Dump(method, args, reply)

	if err != nil {
		log.Println("Failed when calling service for method:", method, err)
	}

	return err
}
