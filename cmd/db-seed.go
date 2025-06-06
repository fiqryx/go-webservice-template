package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"template.go/registry"
)

func NewDBSeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db:seed",
		Short: "Database seeding",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("Starting database seeding...")
			db, err := initDB(cmd)
			if err != nil {
				slog.Error("Error initializing database: ", slog.Any("error", err))
				return
			}

			for _, f := range registry.DBRegistry.GetFactories() {
				if err := f(db); err != nil {
					slog.Error("Error factory", slog.Any("error", err))
					continue
				}
			}

			slog.Info("Database seeding completed!")
		},
	}

	cmd.Flags().IntP("count", "c", 1, "Number of records to seed")
	return cmd
}
