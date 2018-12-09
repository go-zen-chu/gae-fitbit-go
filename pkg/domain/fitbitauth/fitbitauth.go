package fitbitauth

import "net/http"

type FitbitAuthHandler interface {
	HandleFitbitAuth(w http.ResponseWriter, r *http.Request)
}

type fitbitAuthHandler struct{}

func NewFitbitAuthHandler() FitbitAuthHandler {
	return &fitbitAuthHandler{}
}

func (fah *fitbitAuthHandler) HandleFitbitAuth(w http.ResponseWriter, r *http.Request) {

}
