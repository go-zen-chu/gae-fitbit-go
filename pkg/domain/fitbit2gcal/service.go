package fitbit2gcal

import "net/http"

type Service interface {
	HandleFitbit2GCal(w http.ResponseWriter, r *http.Request)
}

type service struct {
	fitbitClient FitbitClient
	gcalClient GCalClient
}

const (
	dateLayout = "20060102"
)

func (s *service) HandleFitbit2GCal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromDateStr := q.Get("fdate")
	toDateStr := q.Get("tdate")
}

// GetFitbitData : Get sleep, activity duration data from fitbit
func GetFitbitData(fromDateStr, toDateStr string) {

}

func ConvertSleep2Schedule() {

}

func ConvertActivity2Schedule() {

}

func Post2GCal() {

}
