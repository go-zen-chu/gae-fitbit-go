package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-zen-chu/gae-fitbit-go/internal/di"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	dateLayout = "20060102"
)

type Handler interface {
	GetIndex(w http.ResponseWriter, r *http.Request)
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	GetFitbitAuthCode(w http.ResponseWriter, r *http.Request)
	Redirect2GCal(w http.ResponseWriter, r *http.Request)
	GetGCalAuthCode(w http.ResponseWriter, r *http.Request)
	GetFitbit2GCal(w http.ResponseWriter, r *http.Request)
	GetFitbit2GCalToday(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	di di.DI
}

func NewHandler(di di.DI) Handler {
	return &handler{
		di: di,
	}
}

func (h *handler) GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Infof("Request to index from : %s, %v", r.RemoteAddr, r)
}

func (h *handler) Redirect2Fitbit(w http.ResponseWriter, r *http.Request) {
	authURL := h.di.FitbitOAuthClient().GetAuthCodeURL()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

func (h *handler) GetFitbitAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		log.Errorln("Could not get auth code from request")
		http.Error(w, "Failed to handle Fitbit Auth Code", 500)
		return
	}
	code := keys[0]
	if err := h.di.AuthService().HandleFitbitAuthCode(code); err != nil {
		log.Errorln(errors.Wrap(err, "Error while storing token"))
		http.Error(w, "Error while handling Fitbit token", 500)
		return
	}
	log.Info("Success handling fitbit token")
	fmt.Fprintf(w, "OK")
}

func (h *handler) Redirect2GCal(w http.ResponseWriter, r *http.Request) {
	authURL := h.di.GCalOAuthClient().GetAuthCodeURL()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

func (h *handler) GetGCalAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		log.Errorln("Could not get auth code from request")
		http.Error(w, "Error getting Google Calendar Auth Code", 500)
		return
	}
	code := keys[0]
	if err := h.di.AuthService().HandleGCalAuthCode(code); err != nil {
		log.Errorln(errors.Wrap(err, "Error while storing token"))
		http.Error(w, "Error while handling GCal token", 500)
		return
	}
	log.Info("Success handling gcal token")
	fmt.Fprintf(w, "OK")
}

func (h *handler) GetFitbit2GCal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromDateStr := q.Get("from_date")
	toDateStr := q.Get("to_date")
	log.Debugf("request %s %s", fromDateStr, toDateStr)

	var err error
	if fromDateStr == "" || toDateStr == "" {
		log.Errorln(errors.New("Insufficient params fromDate, toDate"))
		http.Error(w, "Insufficient params fromDate, toDate", http.StatusBadRequest)
		return
	}
	fromDate, err := time.Parse(dateLayout, fromDateStr)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error parsing from_date"))
		http.Error(w, "Error parsing from_date", http.StatusBadRequest)
		return
	}
	toDate, err := time.Parse(dateLayout, toDateStr)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error parsing to_date"))
		http.Error(w, "Error parsing to_date", http.StatusBadRequest)
		return
	}
	err = h.di.Fitbit2GCalService().Fitbit2GCal(fromDate, toDate)
	if err != nil {
		log.Errorln(err.Error())
		http.Error(w, "Error handling data", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}

func (h *handler) GetFitbit2GCalToday(w http.ResponseWriter, r *http.Request) {
	err := h.di.Fitbit2GCalService().Fitbit2GCalToday()
	if err != nil {
		log.Errorln(err.Error())
		http.Error(w, "Error handling data", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}
