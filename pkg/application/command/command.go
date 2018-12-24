package command

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
	if2g "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbit2gcal"
	ifba "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbitauth"
	iga "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/gcalauth"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	apiVersion = "v1"
)

type Command interface {
	Run() error
}

type config struct {
	port string
}

type command struct{}

func NewCommand() Command {
	return &command{}
}

var (
	port    = kingpin.Flag("port", "Port of application").Default("8080").Envar("GAE_FITBIT_GO_PORT").String()
	verbose = kingpin.Flag("verbose", "Verbosing application").Short('v').Default("false").Bool()
	// needs token of fitbit and gcal
	fbClientID        = kingpin.Flag("fb-client-id", "Fitbit Client ID").Envar("GAE_FITBIT_GO_FITBIT_CLIENT_ID").String()
	fbClientSecret    = kingpin.Flag("fb-client-secret", "Fitbit Client Secret").Envar("GAE_FITBIT_GO_FITBIT_CLIENT_SECRET").String()
	fbAuthRedirectURI = kingpin.Flag("fb-auth-redirect-uri", "Fitbit auth redirect url").Envar("GAE_FITBIT_GO_FITBIT_AUTH_REDIRECT_URI").String()
)

// Run() : runs http api with specified config
func (c *command) Run() error {
	kingpin.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	cnf := &config{
		port: *port,
	}

	fitbitAuthParams := &dfba.FitbitAuthParams{
		ClientID:     *fbClientID,
		Scope:        "sleep activity",
		RedirectURI:  *fbAuthRedirectURI,
		ResponseType: "code",
		Expires:      "2592000", // 1 week
	}

	fitbitTokenParams := &dfba.FitbitTokenParams{
		ClientID:     *fbClientID,
		ClientSecret: *fbClientSecret,
		GrantType:    "authorization_code",
		RedirectURI:  *fbAuthRedirectURI,
	}
	// create handlers
	indexHandler := index.NewIndexHandler()
	fbaFactory := ifba.NewFactory()
	fbaHandler := fbaFactory.FitbitAuthHandler(fitbitAuthParams, fitbitTokenParams)
	f2gFactory := if2g.NewFactory()
	f2gService := f2gFactory.Service()

	credBytes, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(credBytes, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	gaFactory := iga.NewFactory()
	gaHandler := dga.NewGCalAuthHandler(gaFactory, config)

	// Register http handler to routes
	http.HandleFunc("/index.html", indexHandler.HandleIndex)
	http.HandleFunc(fmt.Sprintf("/%s/fitbitauth", apiVersion), fbaHandler.Redirect2Fitbit)
	http.HandleFunc(fmt.Sprintf("/%s/fitbitstoretoken", apiVersion), fbaHandler.HandleFitbitAuthCode)
	http.HandleFunc(fmt.Sprintf("/%s/fitbit2gcal", apiVersion), f2gService.HandleFitbit2GCal)
	http.HandleFunc(fmt.Sprintf("/%s/gcalauth", apiVersion), gaHandler.Redirect2GCal)
	http.HandleFunc(fmt.Sprintf("/%s/gcalstoretoken", apiVersion), gaHandler.HandleGCalAuthCode)

	log.Infof("Running gae-fitbit-go on : %s", cnf.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", cnf.port), nil)
}
