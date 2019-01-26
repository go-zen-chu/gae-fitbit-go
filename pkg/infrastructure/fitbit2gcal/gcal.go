package fitbit2gcal

import (
	"context"
	"errors"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dga"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
	"google.golang.org/api/calendar/v3"
	log "github.com/sirupsen/logrus"
)

type gcalClient struct {
	store      dga.Store
	gcalConfig *df2g.GCalConfig
}

func NewGCalClient(store dga.Store, gcalConfig *df2g.GCalConfig) df2g.GCalClient {
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
		return errors.New("error: unsupported data type")
	}

	token, err := gc.store.FetchGCalToken()
	if err != nil {
		return err
	}
	// make sure to save new token refreshed via oauth2 library
	defer func () {
		if err := gc.store.WriteGCalToken(token); err != nil {
			log.Errorf("%v\n", err)
		}
	}()

	cli := gc.gcalConfig.OauthConfig.Client(context.Background(), token)
	srv, err := calendar.New(cli)
	if err != nil {
		return err
	}
	_, err = srv.Events.Insert(calID, event).Do()
	if err != nil {
		log.Errorf("%v\n", err)
		return err
	}
	return nil
}
