package main

import (
	"flag"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/brokermanage"
	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/logger"
)

var senderConfigFile string

func init() {
	flag.StringVar(&senderConfigFile, "config", "/etc/calendar/sender_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(senderConfigFile)
	logg := logger.New(config.Logger.Level, config.Logger.Type)
	logg.Debug("succes load configuration")

	broker, err := brokermanage.NewBroker(config.Brocker.Endpoint, config.Brocker.Queue)
	if err != nil {
		logg.Error("message", err.Error())
		return
	}

	messages, err := broker.ReadMessages()
	if err != nil {
		logg.Error("message", err.Error())
	}

	var ending chan struct{}

	go func() {
		for message := range messages {
			logg.Info("message", "reading event", "event", message.Body)
		}
	}()

	logg.Info("message", "run reading topic")
	<-ending
}
