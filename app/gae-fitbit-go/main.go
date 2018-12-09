package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-zen-chu/gae-fitbit-go/pkg/application/command"
)

func main() {
	cmd := command.NewCommand()
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
