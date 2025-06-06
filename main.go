package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"template.go/cmd"
	c "template.go/packages/common"
)

var description = `Command Line Interface

A unified tool for managing all application components including:
• HTTP Server - Start/Restart the web service
• Database Migrations - Schema version control
• Data Seeding - Populate database with initial data
• Test Factories - Generate fake data for development`

func main() {
	if err := NewCmd().Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func NewCmd() *cobra.Command {
	debug := c.Env("DEBUG") == "true"

	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	var root = &cobra.Command{
		Use:   "CLI",
		Short: "Command Line Interface tool HTTP server, database migrations, seeding, and factories",
		Long:  description,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	root.PersistentFlags().StringP("dsn", "d", c.Env("DATABASE_URL"), "Database connection string")
	root.PersistentFlags().BoolP("debug", "D", debug, "Show line code on Log and SQL queries")

	root.AddCommand(cmd.NewServeCmd())
	root.AddCommand(cmd.NewMigrateCmd())
	root.AddCommand(cmd.NewDBBackupCmd())
	root.AddCommand(cmd.NewDBSeedCmd())
	root.AddCommand(cmd.NewMDBFactoryCmd())
	root.AddCommand(cmd.NewMakeRepo())
	root.AddCommand(cmd.NewMakeServices())
	root.AddCommand(cmd.NewMakeController())

	return root
}
