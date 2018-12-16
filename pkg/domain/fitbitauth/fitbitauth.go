package fitbitauth

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type FitbitAuthHandler interface {
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request)
}

type fitbitAuthHandler struct {
	FitbitAuthParams *FitbitAuthParams
	factory          Factory
}

func NewFitbitAuthHandler(fac *FitbitAuthParams, fbaf Factory) FitbitAuthHandler {
	return &fitbitAuthHandler{
		FitbitAuthParams: fac,
		factory:          fbaf,
	}
}

// Redirect to Fitbit's oauth url
func (fah *fitbitAuthHandler) Redirect2Fitbit(w http.ResponseWriter, r *http.Request) {
	redirectURL := CreateFitbitAuthURL(fah.FitbitAuthParams)
	http.Redirect(w, r, redirectURL.String(), http.StatusSeeOther)
}

// CreateFitbitAuthURL : Generate Fitbit's oauth authorization flow url
func CreateFitbitAuthURL(fac *FitbitAuthParams) *url.URL {
	url := &url.URL{}
	url.Scheme = "https"
	url.Host = "www.fitbit.com"
	url.Path = "/oauth2/authorize"
	q := url.Query()
	q.Set("client_id", fac.ClientID)
	q.Set("redirect_uri", fac.RedirectURI)
	q.Set("scope", fac.Scope)
	q.Set("expires_in", fac.Expires)
	q.Set("response_type", fac.ResponseType)
	url.RawQuery = q.Encode()
	log.Debugf("create fitbit auth url : %s", url.String())
	return url
}

// HandleFitbitAuthCode : Will recieve auth code from Fitbit, store it to Store
func (fah *fitbitAuthHandler) HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		log.Errorf("Could not get auth code from request")
		return
	}

	st, err := fah.factory.FileStore()
	if err != nil {
		log.Errorf("Error while getting store : %v", err)
		return
	}

	err = st.WriteAuthCode(keys[0])
	if err != nil {
		log.Errorf("Error while writing auth code : %v", err)
		return
	}
}
