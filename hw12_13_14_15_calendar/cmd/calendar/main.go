package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/app"
	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.Type)
	logg.Debug("succes load configuration")
	fmt.Println(config.Storage.DB.Endpoint, config.Storage.DB.Database,
		config.Storage.DB.User, config.Storage.DB.Pass)

	var storage app.Storage
	switch config.Storage.Type {
	case "inmemory":
		storage = memorystorage.New()
	case "db":
		stor, err := sqlstorage.New(config.Storage.DB.Endpoint, config.Storage.DB.Database,
			config.Storage.DB.User, config.Storage.DB.Pass)
		if err != nil {
			logg.Error("message", err.Error())
			return
		}
		storage = stor
	}
	calendar := app.New(logg, storage, config.Server.Port, config.Server.Host)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)
	grpcServer := internalgrpc.NewServer(logg, calendar, config.GRPCServer.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("message", "calendar is running...")

	go func() {
		grpcServer.Start(ctx)
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
