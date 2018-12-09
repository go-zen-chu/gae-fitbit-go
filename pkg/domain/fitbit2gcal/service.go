package fitbit2gcal

import "net/http"

type Service interface {
	HandleFitbit2GCal(w http.ResponseWriter, r *http.Request)
}
