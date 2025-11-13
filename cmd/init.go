package cmd

import (
	"fmt"
	"os"

	"github.com/lmorchard/linkding-to-markdown/internal/templates"
	"github.com/spf13/cobra"
)

const defaultConfigContent = `# Configuration file for linkding-to-markdown
# Copy this to linkding-to-markdown.yaml and customize as needed

# Database configuration (optional, currently not used)
database: "linkding-to-markdown.db"

# Logging configuration
verbose: false
debug: false
log_json: false

# Linkding instance configuration
linkding:
  # URL of your Linkding instance (required)
  url: "https://linkding.example.com"

  # API token (generate in Linkding settings)
  # You can also set this via environment variable: LINKDING_TO_MARKDOWN_LINKDING_TOKEN
  token: "your-api-token-here"

  # HTTP request timeout
  timeout: 30s

# Fetch command configuration
fetch:
  # Number of days to fetch (from now, default: 7)
  days: 7

  # Or specify explicit date range (YYYY-MM-DD format)
  # since: "2025-01-01"
  # until: "2025-01-31"

  # Search query to filter bookmarks
  # query: "golang"

  # Output file (leave empty for stdout)
  # output: "bookmarks.md"

  # Document title
  title: "Bookmarks"

  # Output options
  no_notes: false
  no_tags: false
  no_group_by_date: false
  date_format: "2006-01-02"

  # Template file to use for output (leave empty to use built-in default)
  # template: "linkding-to-markdown.md"
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration and template files",
	Long: `Create default configuration file and custom template file for customization.

This command generates:
  - linkding-to-markdown.yaml (configuration file)
  - linkding-to-markdown.md (customizable template, or use --template-file to specify)

Use --force to overwrite existing files.

Example:
  linkding-to-markdown init
  linkding-to-markdown init --template-file my-template.md
  linkding-to-markdown init --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := GetLogger()
		force, _ := cmd.Flags().GetBool("force")
		templateFile, _ := cmd.Flags().GetString("template-file")

		configFile := "linkding-to-markdown.yaml"

		// Check if config file exists
		configExists := fileExists(configFile)
		if configExists && !force {
			return fmt.Errorf("config file %s already exists (use --force to overwrite)", configFile)
		}

		// Check if template file exists
		templateExists := fileExists(templateFile)
		if templateExists && !force {
			return fmt.Errorf("template file %s already exists (use --force to overwrite)", templateFile)
		}

		// Create config file
		if err := os.WriteFile(configFile, []byte(defaultConfigContent), 0o644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}

		if configExists {
			log.Infof("Overwrote %s", configFile)
		} else {
			log.Infof("Created %s", configFile)
		}

		// Get default template content
		templateContent, err := templates.GetDefaultTemplate()
		if err != nil {
			return fmt.Errorf("failed to get default template: %w", err)
		}

		// Create template file
		if err := os.WriteFile(templateFile, []byte(templateContent), 0o644); err != nil {
			return fmt.Errorf("failed to create template file: %w", err)
		}

		if templateExists {
			log.Infof("Overwrote %s", templateFile)
		} else {
			log.Infof("Created %s", templateFile)
		}

		fmt.Printf("\n✅ Initialization complete!\n\n")
		fmt.Printf("Next steps:\n")
		fmt.Printf("  1. Edit %s and add your Linkding URL and API token\n", configFile)
		fmt.Printf("  2. (Optional) Customize %s for your preferred output format\n", templateFile)
		fmt.Printf("  3. Run: linkding-to-markdown fetch --output bookmarks.md\n\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("force", false, "Overwrite existing files")
	initCmd.Flags().String("template-file", "linkding-to-markdown.md", "Name of custom template file to create")
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
