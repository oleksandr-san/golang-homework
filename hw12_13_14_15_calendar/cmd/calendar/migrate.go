package main

import (
	"fmt"
	"log"

	_ "github.com/alexandera5/hw12_13_14_15_calendar/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate <command> [<args>]",
	Short: "Migrate calendar database",
	Long:  `This command migrates calendar database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		return migrate(configPath, args[0], args[1:]...)
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().String("config", "", "configuration file path")
	migrateCmd.MarkFlagRequired("config")
}

func migrate(configPath string, command string, args ...string) error {
	config, err := ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("migrate: failed to open DB: %w", err)
	}

	db, err := goose.OpenDBWithDriver(config.Storage.Driver, config.Storage.Dsn)
	if err != nil {
		return fmt.Errorf("migrate: failed to open DB: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("migrate: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.Run(command, db, ".", args...); err != nil {
		return fmt.Errorf("migrate: %v: %w", command, err)
	}

	return nil
}
