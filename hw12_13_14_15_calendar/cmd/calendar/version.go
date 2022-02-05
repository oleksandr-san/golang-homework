package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get application version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := json.NewEncoder(os.Stdout).Encode(struct {
			Release   string
			BuildDate string
			GitHash   string
		}{
			Release:   release,
			BuildDate: buildDate,
			GitHash:   gitHash,
		}); err != nil {
			return fmt.Errorf("error while decode version info: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
