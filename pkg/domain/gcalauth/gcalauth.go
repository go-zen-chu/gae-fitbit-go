//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package gcalauth

import (
	"context"
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
	config  *oauth2.Config
}

func NewGCalAuthHandler(gaf Factory, config *oauth2.Config) GCalAuthHandler {
	return &gcalAuthHandler{
		factory: gaf,
		config:  config,
	}
}

// Redirect to GCal's oauth url
func (gah *gcalAuthHandler) Redirect2GCal(w http.ResponseWriter, r *http.Request) {
	authURL := gah.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
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

	fst, _ := gah.factory.FileStore()
	// if err != nil {
	// 	err = errors.Wrap(err, "Error while getting store")
	// 	log.Errorln(err)
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	token, err := gah.config.Exchange(context.TODO(), code)
	if err != nil {
		err = errors.Wrap(err, "Error while getting token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = fst.WriteGCalTokens(token)
	if err != nil {
		err = errors.Wrap(err, "Error while storing token")
		log.Errorln(err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Info("Success storing gcal tokens")
	fmt.Fprintf(w, "OK")
}
