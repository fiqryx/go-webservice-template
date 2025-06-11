package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewGenerateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate:key",
		Short: "Generate secret key",
		Run: func(cmd *cobra.Command, args []string) {
			key := make([]byte, 32)
			_, err := rand.Read(key)
			if err != nil {
				fmt.Println("Error generating key:", err)
				return
			}

			envPath := ".env"
			envKey := "KEY="
			encodedKey := fmt.Sprintf(`"%s"`, base64.URLEncoding.EncodeToString(key))

			file, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Println("Error opening .env file:", err)
				return
			}
			defer file.Close()

			content, err := os.ReadFile(envPath)
			if err != nil {
				fmt.Println("Error reading .env file:", err)
				return
			}

			var newContent string
			if contains(string(content), envKey) {
				newContent = writeEnv(string(content), envKey, envKey+encodedKey)
			} else {
				newContent = string(content) + "\n" + envKey + encodedKey + "\n"
			}

			err = os.WriteFile(envPath, []byte(newContent), 0644)
			if err != nil {
				fmt.Println("Error writing to .env file:", err)
				return
			}
			fmt.Println("Key generated")
		},
	}

	return cmd
}

func contains(content string, prefix string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

// write to replace or add a line
func writeEnv(content string, prefix string, newLine string) string {
	lines := strings.Split(content, "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = newLine
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, newLine)
	}
	return strings.Join(lines, "\n")
}
