package index

import (
	"fmt"
	"net/http"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
)

type indexHandler struct{}

// NewIndexHandler : Get concrete IndexHandler
func NewIndexHandler() index.IndexHandler {
	return &indexHandler{}
}

func (ih *indexHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
