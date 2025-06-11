package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"webservices/registry"
)

func NewDBBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db:backup",
		Short: "Run database backup",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("Starting database backup...")
			db, err := initDB(cmd)
			if err != nil {
				slog.Error("Error initializing database: ", slog.Any("error", err))
				return
			}

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				slog.Error("Error getting output directory", slog.Any("error", err))
				return
			}

			if err := registry.Database.Backup(db, output); err != nil {
				slog.Error("Backup failed", slog.Any("error", err))
			}
		},
	}

	output := fmt.Sprintf("./storage/backup/%s", time.Now().Format("20060102"))
	cmd.Flags().StringP("output", "o", output, "Output directory for backup files")

	return cmd
}
