package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-zhy/secm/pkg/errors"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/open-zhy/secm/pkg/secret"
	"github.com/open-zhy/secm/pkg/workspace"
	"github.com/spf13/cobra"
)

var (
	showTags bool
	showDesc bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored secrets",
	Long:  `Display a list of all stored secrets with their IDs and creation times.`,
	RunE:  runList,
}

func init() {
	listCmd.Flags().BoolVarP(&showTags, "tags", "t", false, "Show secret tags")
	listCmd.Flags().BoolVarP(&showDesc, "description", "d", false, "Show secret descriptions")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Load workspace
	ws, err := workspace.Load(profile)
	if err != nil {
		return errors.Wrapf(err, "failed to load workspace")
	}

	// Read secrets directory
	entries, err := os.ReadDir(ws.SecretsDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read secrets directory")
	}

	if len(entries) == 0 {
		screen.Println("No secrets found")
		return nil
	}

	// Print header
	headers := []string{"ID", "Name", "Type"}
	if showDesc {
		headers = append(headers, "Description")
	}
	if showTags {
		headers = append(headers, "Tags")
	}
	headers = append(headers, "Created At")

	// Calculate column widths
	widths := map[string]int{
		"ID":          36,
		"Name":        30,
		"Type":        15,
		"Description": 30,
		"Tags":        30,
		"Created At":  20,
	}

	// Print headers
	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds", widths["ID"], widths["Name"], widths["Type"])
	if showDesc {
		format += fmt.Sprintf("  %%-%ds", widths["Description"])
	}
	if showTags {
		format += fmt.Sprintf("  %%-%ds", widths["Tags"])
	}
	format += fmt.Sprintf("  %%-%ds\n", widths["Created At"])

	screen.Println("Available secrets:")
	headerInterface := make([]interface{}, len(headers))
	for i, v := range headers {
		headerInterface[i] = v
	}
	screen.Printf(format, headerInterface...)
	screen.Println(strings.Repeat("-", calculateLineWidth(widths, showDesc, showTags)))

	// List secrets
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		secretPath := filepath.Join(ws.SecretsDir, entry.Name())
		s, err := secret.Load(secretPath)
		if err != nil {
			continue
		}

		// Prepare values
		id := strings.TrimSuffix(entry.Name(), ".yml")
		values := []interface{}{
			truncate(id, widths["ID"]),
			truncate(s.Name, widths["Name"]),
			truncate(s.Type, widths["Type"]),
		}
		if showDesc {
			values = append(values, truncate(s.Description, widths["Description"]))
		}
		if showTags {
			values = append(values, truncate(strings.Join(s.Tags, ", "), widths["Tags"]))
		}
		values = append(values, s.CreatedAt.Format("2006-01-02 15:04:05"))

		screen.Printf(format, values...)
	}

	return nil
}

func calculateLineWidth(widths map[string]int, showDesc, showTags bool) int {
	width := widths["ID"] + widths["Name"] + widths["Type"] + widths["Created At"] + 8 // 8 for spacing
	if showDesc {
		width += widths["Description"] + 2
	}
	if showTags {
		width += widths["Tags"] + 2
	}
	return width
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
