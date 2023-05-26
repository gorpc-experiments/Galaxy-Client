package GalaxyClient

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"log"
	"os"
	"strconv"
)

func GetRPCPort() int {
	port := 0
	if KubernetesUtil.IsRunningInKubernetes() {
		port = KubernetesUtil.GetInternalServicePort()
	}
	if port == 0 {
		envPort := os.Getenv("RPC_PORT")
		if envPort == "" {
			log.Fatalln("RPC_PORT env variable isn't set")
		}
		envport, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatalln(err.Error())
		}
		port = envport
	}

	return port
}
