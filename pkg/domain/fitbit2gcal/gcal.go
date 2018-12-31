package fitbit2gcal

import "google.golang.org/api/calendar/v3"

type GCalClient interface {
	// PostSchedule() error
	InsertEvent(event *calendar.Event, dataType string) error
}
