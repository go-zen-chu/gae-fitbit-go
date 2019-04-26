//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package gcalauth

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// FitbitAuthHandler : Manage any authorization process against Fitbit API
type GCalAuthHandler interface {
	Redirect2GCal(w http.ResponseWriter, r *http.Request)
	HandleGCalAuthCode(w http.ResponseWriter, r *http.Request)
}

type gcalAuthHandler struct {
	factory     Factory
	store       Store
	oauthConfig *oauth2.Config
	oauthClient OAuthClient
}

func NewGCalAuthHandler(gaf Factory, store Store, oauthConfig *oauth2.Config) GCalAuthHandler {
	oauthClient := gaf.OAuthClient(oauthConfig)
	return &gcalAuthHandler{
		factory:     gaf,
		store:       store,
		oauthConfig: oauthConfig,
		oauthClient: oauthClient,
	}
}

// Redirect to GCal's oauth url
func (gah *gcalAuthHandler) Redirect2GCal(w http.ResponseWriter, r *http.Request) {
	authURL := gah.oauthClient.GetAuthCodeURL()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// HandleFitbitAuthCode : Will recieve auth code from Fitbit, store it
func (gah *gcalAuthHandler) HandleGCalAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	var err error
	if !ok || len(keys[0]) < 1 {
		log.Errorln("Could not get auth code from request")
		http.Error(w, "Error getting Google Calendar Auth Code", 500)
		return
	}
	// auth code is one time, no need to save it
	code := keys[0]
	log.Debugf("auth code :%s", code)

	token, err := gah.oauthClient.Exchange(code)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error while getting token"))
		http.Error(w, "Error getting Google Calendar token", 500)
		return
	}

	err = gah.store.WriteGCalToken(token)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error while storing token"))
		http.Error(w, "Error storing Google Calendar token", 500)
		return
	}
	log.Info("Success storing gcal tokens")
	fmt.Fprintf(w, "OK")
}
