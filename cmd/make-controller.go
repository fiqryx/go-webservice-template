package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewMakeController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:controller",
		Short: "Create new controllers",
		Run: func(cmd *cobra.Command, args []string) {
			name, err := cmd.Flags().GetString("name")
			if err != nil || name == "" {
				slog.Error("Required flag --name not provided")
				return
			}

			filename := fmt.Sprintf("%s.repository.go", name)
			outputDir, err := cmd.Flags().GetString("output")
			if err != nil {
				slog.Error("Getting output directory: %v", slog.Any("error", err))
				return
			}

			writeSource(name, filepath.Join(outputDir, filename), controllerCode)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Controller name (required)")
	cmd.Flags().StringP("output", "o", "./src/controllers", "Output directory for controller")

	return cmd
}

const controllerCode = `package controllers

type {{.Name}}Controller struct {
	// your services...
}

// inject depedencies on the params, and adjust on [registry/controller.go]
func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{
		// your services...
	}
}
`
