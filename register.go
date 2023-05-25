package GalaxyClient

type RegisterRequest struct {
	Address    string
	Components []string
}

type RegisterResponse struct {
	Success bool
}

func (galaxy *GalaxyClient) RegisterToGalaxy(srvc any) {
	var response RegisterResponse
	err := galaxy.client.Call("Galaxy.Register", RegisterRequest{"127.0.0.1:3456", ExportList(srvc)}, &response)
	if err != nil {
		println("Failed to register service: ", err)
	}

}
