package cmd

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"
	c "webservices/packages/common"
	"webservices/packages/file"
)

func NewMakeServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:service",
		Short: "Create new service",
		Run: func(cmd *cobra.Command, args []string) {
			name, err := cmd.Flags().GetString("name")
			if err != nil || name == "" {
				slog.Error("Required flag --name not provided")
				return
			}

			filename := fmt.Sprintf("%s.service.go", name)
			outputDir, err := cmd.Flags().GetString("output")
			if err != nil {
				slog.Error("Getting output directory: %v", slog.Any("error", err))
				return
			}

			data := map[string]any{
				"Name": c.ToUpper(c.ToCamelCase(name)),
			}

			file.Create(filepath.Join(outputDir, filename), serviceCode, &data)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Services name (required)")
	cmd.Flags().StringP("output", "o", "./src/services", "Output directory for services")

	return cmd
}

const serviceCode = `package services

type {{.Name}}Service struct {
	*baseService
	// other repositories...
}

// inject depedencies on the params, and adjust on [registry/services.go]
func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{
		baseService: newBaseService(),
		// other repositories...
	}
}
`
