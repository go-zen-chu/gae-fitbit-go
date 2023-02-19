package fitbit2gcal

import (
	"context"

	"github.com/pkg/errors"

	da "github.com/go-zen-chu/gae-fitbit-go/domain/auth"
	df2g "github.com/go-zen-chu/gae-fitbit-go/domain/fitbit2gcal"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	calendar "google.golang.org/api/calendar/v3"
)

type gcalClient struct {
	store      da.Store
	gcalConfig *df2g.GCalConfig
}

func NewGCalClient(store da.Store, gcalConfig *df2g.GCalConfig) df2g.GCalClient {
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
	ctx := context.Background()
	tokenSource := gc.gcalConfig.OauthConfig.TokenSource(ctx, token)
	cli := oauth2.NewClient(ctx, tokenSource)
	defer func() {
		err := gc.updateStoredToken(tokenSource)
		if err != nil {
			log.Errorf("%v\n", err)
		}
	}()
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

// make sure to save new token refreshed via oauth2 library
func (gc *gcalClient) updateStoredToken(tokenSource oauth2.TokenSource) error {
	latestToken, err := tokenSource.Token()
	if err != nil {
		return errors.Wrap(err, "Error getting latest token")
	}
	if err = gc.store.WriteGCalToken(latestToken); err != nil {
		return errors.Wrap(err, "Error writing latest token")
	}
	return nil
}
