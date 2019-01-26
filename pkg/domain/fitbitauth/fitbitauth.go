//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

// FitbitAuthHandler : Manage any authorization process against Fitbit API
type FitbitAuthHandler interface {
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request)
}

type fitbitAuthHandler struct {
	factory           Factory
	oauthConfig *oauth2.Config
	oauthClient            OAuthClient
}

func NewFitbitAuthHandler(fbaf Factory, oauthConfig *oauth2.Config) FitbitAuthHandler {
	oauthClient := fbaf.OAuthClient(oauthConfig)
	return &fitbitAuthHandler{
		factory:           fbaf,
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

	token, err := fah.oauthClient.Exchange(code)
	if err != nil {
		err = errors.Wrap(err, "Error while getting token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = fst.WriteFitbitToken(token)
	if err != nil {
		err = errors.Wrap(err, "Error while storing token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Info("Success storing fitbit tokens")
	fmt.Fprintf(w, "OK")
}
