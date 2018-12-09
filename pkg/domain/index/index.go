package index

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// IndexHandler : Only handles index
type IndexHandler interface {
	HandleIndex(w http.ResponseWriter, r *http.Request)
}

type indexHandler struct{}

// NewIndexHandler : Get concrete IndexHandler
func NewIndexHandler() IndexHandler {
	return &indexHandler{}
}

func (ih *indexHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Infof("Request to index from : %s", r.RemoteAddr)
}
