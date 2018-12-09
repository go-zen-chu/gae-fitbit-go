package command

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
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
	port = kingpin.Flag("port", "Port of application").Default("8080").Envar("GAE_FITBIT_GO_PORT").String()
	// needs token of fitbit and gcal
)

// Run() : runs http api with specified config
func (c *command) Run() error {
	kingpin.Parse()
	cnf := &config{
		port: *port,
	}

	ih := index.NewIndexHandler()
	fah := fitbitauth.NewFitbitAuthHandler()

	// TODO: cred を作成して、それを infra の service にわたす

	// Register http handler to routes
	http.HandleFunc("/index.html", ih.HandleIndex)
	http.HandleFunc(fmt.Sprintf("/%s/fitbitauth", apiVersion), fah.HandleFitbitAuth)

	log.Printf("Running gae-fitbit-go on : %s", cnf.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", cnf.port), nil)
}
