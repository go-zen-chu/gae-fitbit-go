//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbit2gcal

import calendar "google.golang.org/api/calendar/v3"

type GCalClient interface {
	// PostSchedule() error
	InsertEvent(event *calendar.Event, dataType string) error
}
