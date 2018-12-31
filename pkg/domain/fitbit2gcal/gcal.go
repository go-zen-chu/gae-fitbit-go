//go:generate mockgen -source gcal.go -destination mock_gcal.go
package fitbit2gcal

import "google.golang.org/api/calendar/v3"

type GCalClient interface {
	// PostSchedule() error
	InsertEvent(event *calendar.Event, dataType string) error
}
