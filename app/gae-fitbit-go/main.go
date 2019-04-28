package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/application/command"
)

func main() {
	httpServer := command.NewHttpServer()
	// run command with actual configuration
	cmd := command.NewCommand(httpServer)

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
