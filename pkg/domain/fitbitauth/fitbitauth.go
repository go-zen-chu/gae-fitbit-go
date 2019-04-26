//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// FitbitAuthHandler : Manage any authorization process against Fitbit API
type FitbitAuthHandler interface {
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request)
}

type fitbitAuthHandler struct {
	factory     Factory
	store       Store
	oauthConfig *oauth2.Config
	oauthClient OAuthClient
}

func NewFitbitAuthHandler(fbaf Factory, store Store, oauthConfig *oauth2.Config) FitbitAuthHandler {
	oauthClient := fbaf.OAuthClient(oauthConfig)
	return &fitbitAuthHandler{
		factory:     fbaf,
		store:       store,
		oauthConfig: oauthConfig,
		oauthClient: oauthClient,
	}
}

// Redirect to Fitbit's oauth url
func (fah *fitbitAuthHandler) Redirect2Fitbit(w http.ResponseWriter, r *http.Request) {
	authURL := fah.oauthClient.GetAuthCodeURL()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// HandleFitbitAuthCode : Will recieve auth code from Fitbit, store it
func (fah *fitbitAuthHandler) HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	var err error
	if !ok || len(keys[0]) < 1 {
		log.Errorln("Could not get auth code from request")
		http.Error(w, "Failed to handle Fitbit Auth Code", 500)
		return
	}
	// auth code is one time, no need to save it
	code := keys[0]

	token, err := fah.oauthClient.Exchange(code)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error while getting token"))
		http.Error(w, "Failed while getting Fitbit token", 500)
		return
	}

	err = fah.store.WriteFitbitToken(token)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error while storing token"))
		http.Error(w, "Error while storing Fitbit token", 500)
		return
	}
	log.Info("Success storing fitbit tokens")
	fmt.Fprintf(w, "OK")
}
