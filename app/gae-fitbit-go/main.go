package main

import (
	"github.com/go-zen-chu/gae-fitbit-go/pkg/domain/index"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbit2gcal"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbitauth"
	"github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/gcalauth"
	log "github.com/sirupsen/logrus"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/application/command"
)

func main() {
	indexHandler := index.NewIndexHandler()
	fitbitauthFactory := fitbitauth.NewFactory()
	gcalauthFactory := gcalauth.NewFactory()
	fitbit2gcalFactory := fitbit2gcal.NewFactory()
	httpServer := command.NewHttpServer()

	// run command with actual configuration
	cmd := command.NewCommand(
		indexHandler,
		fitbitauthFactory,
		gcalauthFactory,
		fitbit2gcalFactory,
		httpServer)

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
