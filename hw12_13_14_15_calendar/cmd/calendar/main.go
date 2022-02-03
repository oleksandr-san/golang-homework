package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandera5/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexandera5/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/alexandera5/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/alexandera5/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func makeStorage(config StorageConf) app.Storage {
	switch config.Driver {
	case "memory":
		return memorystorage.New()
	case "postgresql":
		return sqlstorage.New(config.Driver, config.Dsn)
	default:
		log.Fatal("Unknown storage type")
		return nil
	}
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := ReadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(config)
	logg, err := logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Close()

	storage := makeStorage(config.Storage)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(
		config.Server.Address(), logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
