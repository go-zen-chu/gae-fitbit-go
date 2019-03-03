package command

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
	"google.golang.org/api/calendar/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	apiVersion = "v1"
)

type Command interface {
	Run() error
}

type command struct {
	indexHandler       index.IndexHandler
	fitbitauthFactory  dfba.Factory
	gcalauthFactory    dga.Factory
	fitbit2gcalFactory df2g.Factory
	httpServer         HttpServer
}

func NewCommand(
	indexHandler index.IndexHandler,
	fitbitauthFactory dfba.Factory,
	gcalauthFactory dga.Factory,
	fitbit2gcalFactory df2g.Factory,
	httpServer HttpServer,
) Command {
	return &command{
		indexHandler:       indexHandler,
		fitbitauthFactory:  fitbitauthFactory,
		gcalauthFactory:    gcalauthFactory,
		fitbit2gcalFactory: fitbit2gcalFactory,
		httpServer:         httpServer,
	}
}

var (
	port    = kingpin.Flag("port", "Port of application").Default("8080").Envar("PORT").String()
	verbose = kingpin.Flag("verbose", "Verbosing application").Short('v').Default("false").Bool()
	// fitbit options
	fbClientID        = kingpin.Flag("fb-client-id", "Fitbit Client ID").Envar("GAE_FITBIT_GO_FITBIT_CLIENT_ID").String()
	fbClientSecret    = kingpin.Flag("fb-client-secret", "Fitbit Client Secret").Envar("GAE_FITBIT_GO_FITBIT_CLIENT_SECRET").String()
	fbAuthRedirectURI = kingpin.Flag("fb-auth-redirect-uri", "Fitbit auth redirect url").Default("http://127.0.0.1:8080/v1/fitbitstoretoken").Envar("GAE_FITBIT_GO_FITBIT_AUTH_REDIRECT_URI").String()
	// gcal options
	gcalSleepCalendarID    = kingpin.Flag("gcal-sleep-cal-id", "Google sleep calendar ID").Envar("GAE_FITBIT_GO_FITBIT_GCAL_SLEEP_CAL_ID").String()
	gcalActivityCalendarID = kingpin.Flag("gcal-activity-cal-id", "Google activity calendar ID").Envar("GAE_FITBIT_GO_FITBIT_GCAL_ACTIVITY_CAL_ID").String()
	gcalClientID           = kingpin.Flag("gcal-client-id", "Google Calendar Client ID").Envar("GAE_FITBIT_GO_GCAL_CLIENT_ID").String()
	gcalClientSecret       = kingpin.Flag("gcal-client-secret", "Google Calendar Client Secret").Envar("GAE_FITBIT_GO_GCAL_CLIENT_SECRET").String()
	gcalAuthRedirectURI    = kingpin.Flag("gcal-auth-redirect-uri", "GCal auth redirect url").Default("http://localhost:8080/v1/gcalstoretoken").Envar("GAE_FITBIT_GO_GCAL_AUTH_REDIRECT_URI").String()
	// application options
	useCloudStorage        = kingpin.Flag("use-cloud-storage", "Use Cloud Storage or not. If you deploy as GAE, needs to be true").Envar("USE_CLOUD_STORAGE").Bool()
	cloudStorageBucketName = kingpin.Flag("cloud-storage-bucket-name", "Google Cloud bucket name to use").Envar("CLOUD_STORAGE_BUCKET_NAME").String()
)

// Run() : runs http api with specified config
func (c *command) Run() error {
	kingpin.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	var err error
	var fbaStore dfba.Store
	var gaStore dga.Store
	if *useCloudStorage {
		fbaStore, err = c.fitbitauthFactory.CloudStorageStore(*cloudStorageBucketName)
		if err != nil {
			return errors.Wrap(err, "Error while making cloud storage store")
		}
		gaStore, err = c.gcalauthFactory.CloudStorageStore(*cloudStorageBucketName)
		if err != nil {
			return errors.Wrap(err, "Error while making cloud storage store")
		}
	} else {
		fbaStore, err = c.fitbitauthFactory.FileStore()
		if err != nil {
			return errors.Wrap(err, "Error while making file store")
		}
		gaStore, err = c.gcalauthFactory.FileStore()
		if err != nil {
			return errors.Wrap(err, "Error while making file store")
		}
	}

	// create handlers
	fbOauthConfig := &oauth2.Config{
		ClientID:     *fbClientID,
		ClientSecret: *fbClientSecret,
		Endpoint:     fitbit.Endpoint,
		RedirectURL:  *fbAuthRedirectURI,
		Scopes:       []string{"sleep", "activity"},
	}
	fbaHandler := c.fitbitauthFactory.FitbitAuthHandler(fbaStore, fbOauthConfig)

	gcalOauthConfig := &oauth2.Config{
		ClientID:     *gcalClientID,
		ClientSecret: *gcalClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v3/token",
		},
		RedirectURL: *gcalAuthRedirectURI,
		Scopes:      []string{calendar.CalendarEventsScope},
	}
	gaHandler := c.gcalauthFactory.GCalAuthHandler(gaStore, gcalOauthConfig)

	fitbitConfig := &df2g.FitbitConfig{
		OauthConfig: fbOauthConfig,
	}
	gcalConfig := &df2g.GCalConfig{
		SleepCalendarID:    *gcalSleepCalendarID,
		ActivityCalendarID: *gcalActivityCalendarID,
		OauthConfig:        gcalOauthConfig,
	}
	f2gService := c.fitbit2gcalFactory.Service(fitbitConfig, gcalConfig, fbaStore, gaStore)

	// Register http handler to routes
	c.httpServer.HandleFunc("/index.html", c.indexHandler.HandleIndex)

	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbitauth", apiVersion), fbaHandler.Redirect2Fitbit)
	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbitstoretoken", apiVersion), fbaHandler.HandleFitbitAuthCode)

	c.httpServer.HandleFunc(fmt.Sprintf("/%s/gcalauth", apiVersion), gaHandler.Redirect2GCal)
	c.httpServer.HandleFunc(fmt.Sprintf("/%s/gcalstoretoken", apiVersion), gaHandler.HandleGCalAuthCode)

	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbit2gcal", apiVersion), f2gService.HandleFitbit2GCal)
	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbit2gcal/today", apiVersion), f2gService.HandleFitbit2GCal)

	log.Infof("Running gae-fitbit-go on : %s", *port)
	return c.httpServer.Run(*port)
}
