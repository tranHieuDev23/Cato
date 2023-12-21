package rpc

type EchoRequest struct {
	Message string
}

type EchoResponse struct {
	Message string
}

//go:generate genpjrpc -search.name=API -print.place.path_swagger_file=../../../../api/swagger.json
type API interface {
	Echo(EchoRequest) EchoResponse
}
