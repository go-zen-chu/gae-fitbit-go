package command

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
	calendar "google.golang.org/api/calendar/v3"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/di"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/interface/handler"
)

const (
	apiVersion = "v1"
)

type Command interface {
	Run() error
}

type command struct {
	httpServer HttpServer
}

func NewCommand(httpServer HttpServer) Command {
	return &command{
		httpServer: httpServer,
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
	// initialize di and objects
	di := di.NewDI()
	if *useCloudStorage {
		di.InitAuthCloudStorageStore(*cloudStorageBucketName)
	} else {
		di.InitAuthFileStore()
	}
	di.InitFitbitOAuthConfig(&oauth2.Config{
		ClientID:     *fbClientID,
		ClientSecret: *fbClientSecret,
		Endpoint:     fitbit.Endpoint,
		RedirectURL:  *fbAuthRedirectURI,
		Scopes:       []string{"sleep", "activity"},
	})
	di.InitGCalOAuthConfig(&oauth2.Config{
		ClientID:     *gcalClientID,
		ClientSecret: *gcalClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v3/token",
		},
		RedirectURL: *gcalAuthRedirectURI,
		Scopes:      []string{calendar.CalendarEventsScope},
	})
	di.InitGCalConfig(*gcalSleepCalendarID, *gcalActivityCalendarID)

	h := handler.NewHandler(di)

	// Register http handler to routes
	c.httpServer.HandleFunc("/index.html", h.GetIndex)

	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbitauth", apiVersion), h.Redirect2Fitbit)
	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbitstoretoken", apiVersion), h.GetFitbitAuthCode)

	c.httpServer.HandleFunc(fmt.Sprintf("/%s/gcalauth", apiVersion), h.Redirect2GCal)
	c.httpServer.HandleFunc(fmt.Sprintf("/%s/gcalstoretoken", apiVersion), h.GetGCalAuthCode)

	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbit2gcal", apiVersion), h.GetFitbit2GCal)
	c.httpServer.HandleFunc(fmt.Sprintf("/%s/fitbit2gcal/today", apiVersion), h.GetFitbit2GCalToday)

	log.Infof("Running gae-fitbit-go on : %s", *port)
	return c.httpServer.Run(*port)
}
