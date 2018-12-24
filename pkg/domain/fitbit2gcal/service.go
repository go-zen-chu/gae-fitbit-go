package fitbit2gcal

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type Service interface {
	HandleFitbit2GCal(w http.ResponseWriter, r *http.Request)
}

type service struct {
	fitbitClient FitbitClient
	gcalClient   GCalClient
}

const (
	dateLayout = "20060102"
)

func NewService(fbc FitbitClient, gc GCalClient) Service {
	return &service{
		fitbitClient: fbc,
		gcalClient:   gc,
	}
}

func (s *service) HandleFitbit2GCal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromDateStr := q.Get("from_date")
	toDateStr := q.Get("to_date")
	log.Debugf("request %s %s", fromDateStr, toDateStr)

	var err error
	if fromDateStr == "" || toDateStr == "" {
		err = errors.New("Insufficient params fromDate, toDate")
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fromDate, err := time.Parse("20060102", fromDateStr)
	if err != nil {
		err = errors.Wrap(err, "Error parsing fromDateStr")
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	toDate, err := time.Parse("20060102", toDateStr)
	if err != nil {
		err = errors.Wrap(err, "Error parsing toDateStr")
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if fromDate.After(toDate) {
		err = errors.New("Invalid parameter, fromDate > toDate")
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sleeps, activities, err := s.getFitbitData(fromDate, toDate)
	if err != nil {
		err = errors.Wrap(err, "Error getting data from Fitbit")
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("%v %v", sleeps, activities)

	for sleep := range sleeps {
		events :=
		err = s.gcalClient.InsertEvent(event, "sleep")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	for sleep := range sleeps {
		events :=
		err = s.gcalClient.InsertEvent(event, "activity")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	fmt.Fprint(w, "OK")
}

// getFitbitData : Get sleep, activity duration data from fitbit
func (s *service) getFitbitData(fromDate, toDate time.Time) ([]Sleep, []Activity, error) {
	sleeps := make([]Sleep, 0, 5)
	activities := make([]Activity, 0, 5)

	for dt := fromDate; dt.Before(toDate); dt = dt.AddDate(0, 0, 1) {
		dateStr := dt.Format("2006-01-02")
		log.Infof("Getting data of %s", dateStr)
		sd, err := s.fitbitClient.GetSleepData(dateStr)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error while getting sleep data")
		}
		ad, err := s.fitbitClient.GetActivityData(dateStr)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error while getting activity data")
		}
		sleeps = append(sleeps, *sd)
		activities = append(activities, *ad)
	}
	return sleeps, activities, nil
}

func convertSleep2Events(sleep *Sleep) []calendar.Event {
	var 
	sleep.
	return &calendar.Event {
		Summary:
	}
}

func convertActivity2Events() {

}

func Post2GCal() {

}
