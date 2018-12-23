package fitbit2gcal

import "net/http"

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

func (s *service) HandleFitbit2GCal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	_ := q.Get("from_date")
	_ := q.Get("to_date")
}

// getFitbitData : Get sleep, activity duration data from fitbit
func getFitbitData(fromDateStr, toDateStr string) {

}

func ConvertSleep2Schedule() {

}

func ConvertActivity2Schedule() {

}

func Post2GCal() {

}
