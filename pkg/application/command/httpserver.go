//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package command

import (
	"fmt"
	"net/http"
)

type HttpServer interface {
	Run(port string) error
}

type httpServer struct {}

func NewHttpServer() HttpServer {
	return &httpServer{}
}

func (hs *httpServer) Run(port string) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
