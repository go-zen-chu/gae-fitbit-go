//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// FitbitAuthHandler : Manage any authorization process against Fitbit API
type FitbitAuthHandler interface {
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request)
}

type fitbitAuthHandler struct {
	factory           Factory
	fitbitAuthParams  *FitbitAuthParams
	fitbitTokenParams *FitbitTokenParams
}

func NewFitbitAuthHandler(fbaf Factory, fap *FitbitAuthParams, ftp *FitbitTokenParams) FitbitAuthHandler {
	return &fitbitAuthHandler{
		factory:           fbaf,
		fitbitAuthParams:  fap,
		fitbitTokenParams: ftp,
	}
}

// Redirect to Fitbit's oauth url
func (fah *fitbitAuthHandler) Redirect2Fitbit(w http.ResponseWriter, r *http.Request) {
	redirectURL := createFitbitAuthURL(fah.fitbitAuthParams)
	http.Redirect(w, r, redirectURL.String(), http.StatusSeeOther)
}

// createFitbitAuthURL : Generate Fitbit's oauth authorization flow url
func createFitbitAuthURL(fap *FitbitAuthParams) *url.URL {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "www.fitbit.com"
	u.Path = "/oauth2/authorize"
	q := u.Query()
	q.Set("client_id", fap.ClientID)
	q.Set("redirect_uri", fap.RedirectURI)
	q.Set("scope", fap.Scope)
	q.Set("expires_in", fap.Expires)
	q.Set("response_type", fap.ResponseType)
	u.RawQuery = q.Encode()
	log.Debugf("create fitbit auth url : %s", u.String())
	return u
}

// HandleFitbitAuthCode : Will recieve auth code from Fitbit, store it
func (fah *fitbitAuthHandler) HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	var err error
	if !ok || len(keys[0]) < 1 {
		err = errors.New("Could not get auth code from request")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}
	// auth code is one time, no need to save it
	code := keys[0]
	log.Debugf("auth code :%s", code)

	fst, err := fah.factory.FileStore()
	if err != nil {
		err = errors.Wrap(err, "Error while getting store")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	fhc := fah.factory.FitbitHTTPClient()

	ft, err := fhc.GetFitbitToken(fah.fitbitTokenParams, code)
	if err != nil {
		err = errors.Wrap(err, "Error while getting token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = fst.WriteFitbitTokens(ft)
	if err != nil {
		err = errors.Wrap(err, "Error while storing token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Info("Success storing fitbit tokens")
	fmt.Fprintf(w, "OK")
}
