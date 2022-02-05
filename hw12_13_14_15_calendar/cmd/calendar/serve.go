package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandera5/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexandera5/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/alexandera5/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/alexandera5/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve calendar API",
	Long:  `This command starts HTTP API serving.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		return serve(configPath)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("config", "", "configuration file path")
	serveCmd.MarkFlagRequired("config")
}

func connectStorage(config StorageConf, logg *logger.Logger) (app.Storage, context.CancelFunc, error) {
	switch config.Driver {
	case "memory":
		return memorystorage.New(), func() {}, nil
	case "postgresql":
		s := sqlstorage.New(config.Driver, config.Dsn, logg)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		err := s.Connect(ctx)
		if err != nil {
			return nil, cancel, err
		}

		return s, cancel, nil

	default:
		log.Fatal("Unknown storage type")
		return nil, nil, nil
	}
}

func serve(configFile string) error {
	config, err := ReadConfig(configFile)
	if err != nil {
		return err
	}

	log.Println(config)
	logg, err := logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		return err
	}
	defer logg.Close()

	storage, storageCancel, err := connectStorage(config.Storage, logg)
	if err != nil {
		return err
	}
	defer storageCancel()

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
		logg.Error("http serve failed: " + err.Error())
		cancel()
		return err
	}

	return nil
}
