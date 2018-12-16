package command

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
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

	fac := &fitbitauth.FitbitAuthParams{
		ClientID:     *fbClientID,
		Scope:        "sleep activity",
		RedirectURI:  *fbAuthRedirectURI,
		ResponseType: "code",
		Expires:      "2592000", // 1 week
	}

	ih := index.NewIndexHandler()
	fbaf := ifba.NewFactory()
	fah := fitbitauth.NewFitbitAuthHandler(fac, fbaf)

	// TODO: cred を作成して、それを infra の service にわたす

	// Register http handler to routes
	http.HandleFunc("/index.html", ih.HandleIndex)
	http.HandleFunc(fmt.Sprintf("/%s/fitbitauth", apiVersion), fah.Redirect2Fitbit)
	http.HandleFunc(fmt.Sprintf("/%s/fitbitcode", apiVersion), fah.HandleFitbitAuthCode)

	log.Infof("Running gae-fitbit-go on : %s", cnf.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", cnf.port), nil)
}
