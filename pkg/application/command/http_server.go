//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package command

import (
	"fmt"
	"net/http"
)

type HttpServer interface {
	HandleFunc(pattern string, handlerFunc http.HandlerFunc)
	Run(port string) error
}

type httpServer struct {
	mux *http.ServeMux
}

func NewHttpServer() HttpServer {
	return &httpServer{
		mux: http.NewServeMux(),
	}
}

func (hs *httpServer) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	// register to mux so and don't use http.HandleFunc for multiple registration (testing)
	hs.mux.HandleFunc(pattern, handlerFunc)
}

func (hs *httpServer) Run(port string) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", port), hs.mux)
}
