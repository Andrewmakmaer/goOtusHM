package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/brokermanage"
	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/callcearcher"
	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/logger"
	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.Type)
	logg.Debug("succes load configuration")
	storage, err := sqlstorage.New(config.Storage.Endpoint, config.Storage.Database,
		config.Storage.User, config.Storage.Pass)
	if err != nil {
		logg.Error("message", err.Error())
	}

	broker, err := brokermanage.NewBroker(config.Brocker.Endpoint, config.Brocker.Queue)
	if err != nil {
		logg.Error("message", err.Error())
		return
	}

	logg.Info("message", "send event notification")
	eventSearch := scheduler.NewSchedule(5*time.Second, func() {
		callcearcher.SendCallCandidate(storage, &broker, logg)
	})

	eventDelete := scheduler.NewSchedule(5*time.Second, func() {
		callcearcher.DeleteOldEvents(storage, logg)
	})

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	eventSearch.Do(ctx)
	eventDelete.Do(ctx)
	<-ctx.Done()
}
