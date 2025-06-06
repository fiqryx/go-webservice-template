package cmd

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	c "template.go/packages/common"
)

func NewMakeRepo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:repo",
		Short: "Create new repositories",
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

			writeSource(name, filepath.Join(outputDir, filename), repoCode)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repositories name (required)")
	cmd.Flags().StringP("output", "o", "./src/repositories", "Output directory for repositories")

	return cmd
}

func writeSource(name, filePath, source string) {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("create file", slog.Any("error", err))
		return
	}
	defer file.Close()

	tmpl, err := template.New("source").Parse(source)
	if err != nil {
		fmt.Printf("error parse code: %v\n", err)
		return
	}

	data := map[string]any{
		"Name": toUpper(c.ToCamelCase(name)),
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Error write code: %v\n", err)
		return
	}

	fmt.Printf("created: %s\n", filePath)
}

const repoCode = `package repositories
import (
	"gorm.io/gorm"
)

type {{.Name}}Repository struct {
	*baseRepository
}

func New{{.Name}}Repository(db *gorm.DB) *{{.Name}}Repository {
	return &{{.Name}}Repository{
		baseRepository: newBaseRepository(db),
	}
}
`
