package index

import "net/http"

type IndexHandler interface {
	HandleIndex(w http.ResponseWriter, r *http.Request)
}
