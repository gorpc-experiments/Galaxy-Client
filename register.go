package ServiceCore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"os"
	"strconv"
)

type RegisterRequest struct {
	Address    string
	Components []string
}

type RegisterResponse struct {
	Success bool
}

func (galaxy *GalaxyClient) RegisterToGalaxy(srvc any) {
	var response RegisterResponse

	host := KubernetesUtil.GetInternalServiceIP()

	if host == "" {
		host = "127.0.0.1"
	}

	port := KubernetesUtil.GetInternalServicePort()

	if port == 0 {
		portenv := os.Getenv("PORT")

		portval, err := strconv.Atoi(portenv)
		if err != nil {
			println("Failed to parse port: ", err)
		}

		port = portval
	}

	err := galaxy.client.Call("Galaxy.Register", RegisterRequest{fmt.Sprintf("%s:%d", host, port), ExportList(srvc)}, &response)
	if err != nil {
		println("Failed to register service: ", err)
	}

}
