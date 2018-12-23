package command

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
	if2g "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbit2gcal"
	ifba "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbitauth"
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
	fitbitAuthFactory := ifba.NewFactory()
	fitbitAuthHandler := dfba.NewFitbitAuthHandler(fitbitAuthFactory, fitbitAuthParams, fitbitTokenParams)
	fitbit2gcalFactory := if2g.NewFactory()
	fitbit2gcalService := fitbit2gcalFactory.Service()

	// Register http handler to routes
	http.HandleFunc("/index.html", indexHandler.HandleIndex)
	http.HandleFunc(fmt.Sprintf("/%s/fitbitauth", apiVersion), fitbitAuthHandler.Redirect2Fitbit)
	http.HandleFunc(fmt.Sprintf("/%s/storetoken", apiVersion), fitbitAuthHandler.HandleFitbitAuthCode)
	http.HandleFunc(fmt.Sprintf("/%s/fitbit2gcal", apiVersion), fitbit2gcalService.HandleFitbit2GCal)

	log.Infof("Running gae-fitbit-go on : %s", cnf.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", cnf.port), nil)
}
