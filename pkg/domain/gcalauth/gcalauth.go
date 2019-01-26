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
	factory Factory
	oauthConfig *oauth2.Config
	oauthClient  OAuthClient
}

func NewGCalAuthHandler(gaf Factory, oauthConfig *oauth2.Config) GCalAuthHandler {
	oauthClient := gaf.OAuthClient(oauthConfig)
	return &gcalAuthHandler{
		factory: gaf,
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
		err = errors.New("Could not get auth code from request")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}
	// auth code is one time, no need to save it
	code := keys[0]
	log.Debugf("auth code :%s", code)

	fst, err := gah.factory.FileStore()
	if err != nil {
		err = errors.Wrap(err, "Error while getting store")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	token, err := gah.oauthClient.Exchange(code)
	if err != nil {
		err = errors.Wrap(err, "Error while getting token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = fst.WriteGCalToken(token)
	if err != nil {
		err = errors.Wrap(err, "Error while storing token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Info("Success storing gcal tokens")
	fmt.Fprintf(w, "OK")
}
