package fitbit2gcal

import (
	"context"
	"errors"

	"ghe.corp.yahoo.co.jp/pivotal-cf/cluster-spec/config"
	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	calendar "google.golang.org/api/calendar/v3"
)

type gcalClient struct {
	store      df2g.Store
	gcalConfig *df2g.GCalConfig
}

func NewGCalClient(store df2g.Store, gcalConfig *df2g.GCalConfig) df2g.GCalClient {
	return &gcalClient{
		store:      store,
		gcalConfig: gcalConfig,
	}
}

func (gc *gcalClient) InsertEvent(event *calendar.Event, dataType string) error {
	var calID string
	switch dataType {
	case "sleep":
		calID = gc.gcalConfig.SleepCalendarID
	case "activity":
		calID = gc.gcalConfig.ActivityCalendarID
	default:
		return errors.New("Unsupported data type")
	}

	token, err := gc.store.FetchGCalTokens()
	if err != nil {
		return err
	}
	cli := config.Client(context.Background(), token)
	srv, err := calendar.New(cli)
	if err != nil {
		return err
	}
	_, err = srv.Events.Insert(calID, event).Do()
	if err != nil {
		return err
	}
	return nil
}
