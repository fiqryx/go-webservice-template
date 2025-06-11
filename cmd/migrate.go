package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"template.go/database"
	"template.go/registry"
)

func NewMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migration",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("Starting database migration...")
			db, err := initDB(cmd)
			if err != nil {
				slog.Error("Initializing database: ", slog.Any("error", err))
				return
			}

			fresh, _ := cmd.Flags().GetBool("fresh")
			if err := registry.Database.Migrate(db, fresh); err != nil {
				slog.Error("Migration failed", slog.Any("error", err))
			}
		},
	}

	cmd.Flags().BoolP("fresh", "f", false, "force fresh migration")

	return cmd
}

func initDB(cmd *cobra.Command) (*gorm.DB, error) {
	dsn, err := cmd.Flags().GetString("dsn")
	if err != nil {
		return nil, err
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return nil, err
	}

	database.Connect(dsn, debug)
	return database.DB(), nil
}
