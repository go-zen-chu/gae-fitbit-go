package command

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/index"
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
	port = kingpin.Arg("--port", "Port of application").Required().Envar("GAE_FITBIT_GO_PORT").String()
	// needs token of fitbit and gcal
)

// Run() : runs http api with specified config
func (c *command) Run() error {
	cnf := &config{
		port: *port,
	}

	ih := index.NewIndexHandler()

	// TODO: cred を作成して、それを infra の service にわたす

	http.HandleFunc(fmt.Sprintf("/index.html", apiVersion), ih.HandleIndex)

	log.Printf("Running gae-fitbit-go on : %s", cnf.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", cnf.port), nil)
}
