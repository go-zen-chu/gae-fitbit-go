//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbit2gcal

import (
	"fmt"
	"google.golang.org/api/calendar/v3"
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
	fromDate, err := time.Parse(dateLayout, fromDateStr)
	if err != nil {
		err = errors.Wrap(err, "Error parsing fromDateStr")
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	toDate, err := time.Parse(dateLayout, toDateStr)
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

	for _, sleep := range sleeps {
		events, err := convertSleep2Events(&sleep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, event := range events {
			err = s.gcalClient.InsertEvent(&event, "sleep")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
	for _, act := range activities {
		events, err := convertActivity2Events(&act)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, event := range events {
			err = s.gcalClient.InsertEvent(&event, "activity")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
	fmt.Fprint(w, "OK")
}

// getFitbitData : Get sleep, activity duration data from fitbit
func (s *service) getFitbitData(fromDate, toDate time.Time) ([]Sleep, []Activity, error) {
	var sleeps []Sleep
	var activities []Activity

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

func duration2HourMin(duration time.Duration) (time.Duration, time.Duration) {
	durationInMinutes := duration.Round(time.Minute)
	durationHour := durationInMinutes / time.Hour
	var remainDuration time.Duration
	remainDuration = duration - durationHour * time.Hour
	durationMin := remainDuration / time.Minute
	return durationHour, durationMin
}

func convertSleep2Events(sleep *Sleep) ([]calendar.Event, error) {
	var events []calendar.Event

	tz := "Asia/Tokyo"
	lc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	for _, s := range sleep.Sleep{
		duration :=  time.Duration(s.Duration) * time.Millisecond

		startTime, err := time.ParseInLocation("2006-01-02T15:04:05.999", s.StartTime, lc)
		if err != nil {
			return nil, err
		}
		endTime := startTime.Add(duration)

		durationHour, durationMin := duration2HourMin(duration)
		durationAsleep := time.Duration(s.MinutesAsleep) * time.Minute
		asleepHour, asleepMin := duration2HourMin(durationAsleep)
		durationAwake := time.Duration(s.MinutesAwake) * time.Minute
		awakeHour, awakeMin := duration2HourMin(durationAwake)

		summary := fmt.Sprintf("Sleep %02d:%02d", durationHour, durationMin)
		desc := fmt.Sprintf("DateOfSleep : %s\n" +
			"Duration : %02d:%02d\n" +
			"Efficiency : %d\n" +
			"IsMainSleep : %t\n" +
			"MinutesAsleep : %02d:%02d\n" +
			"MinutesAwake : %02d:%02d\n" +
			"LogID : %d\n",
			s.DateOfSleep,
			durationHour, durationMin,
			s.Efficiency,
			s.IsMainSleep,
			asleepHour, asleepMin,
			awakeHour, awakeMin,
			s.LogID)
		ev :=  &calendar.Event {
			Summary: summary,
			Start: &calendar.EventDateTime{
				DateTime: startTime.Format("2006-01-02T15:04:05"),
				TimeZone: tz,
			},
			End: &calendar.EventDateTime{
				DateTime: endTime.Format("2006-01-02T15:04:05"),
				TimeZone: tz,
			},
			Description: desc,
		}
		events = append(events, *ev)
	}
	return events, nil
}

func convertActivity2Events(activity *Activity) ([]calendar.Event, error) {
	var events []calendar.Event

	tz := "Asia/Tokyo"
	lc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	for _, a := range activity.Activities {
		duration :=  time.Duration(a.Duration) * time.Millisecond

		startTime, err := time.ParseInLocation("2006-01-02T15:04:05.999-07:00", a.StartTime, lc)
		if err != nil {
			return nil, err
		}
		endTime := startTime.Add(duration)

		durationHour, durationMin := duration2HourMin(duration)

		summary := fmt.Sprintf("%s %02d:%02d", a.ActivityName, durationHour, durationMin)
		desc := fmt.Sprintf("Calories : %d\n" +
			"Duration : %02d:%02d\n" +
			"Steps : %d\n" +
			"LogID : %d\n",
			a.Calories,
			durationHour, durationMin,
			a.Steps,
			a.LogID)
		ev :=  &calendar.Event {
			Summary: summary,
			Start: &calendar.EventDateTime{
				DateTime: startTime.Format("2006-01-02T15:04:05"),
				TimeZone: tz,
			},
			End: &calendar.EventDateTime{
				DateTime: endTime.Format("2006-01-02T15:04:05"),
				TimeZone: tz,
			},
			Description: desc,
		}
		events = append(events, *ev)
	}
	return events, nil
}

