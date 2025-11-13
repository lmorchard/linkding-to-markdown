package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/lmorchard/linkding-to-markdown/internal/linkding"
	"github.com/lmorchard/linkding-to-markdown/internal/markdown"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch bookmarks from Linkding and generate markdown",
	Long: `Fetch bookmarks from Linkding over a given timespan and produce a markdown file.

Examples:
  # Fetch bookmarks from the last 7 days
  linkding-to-markdown fetch --days 7

  # Fetch bookmarks from the last week with custom output
  linkding-to-markdown fetch --days 7 --output bookmarks.md

  # Fetch bookmarks from a specific date range
  linkding-to-markdown fetch --since 2025-01-01 --until 2025-01-31

  # Fetch bookmarks with a search query
  linkding-to-markdown fetch --query "golang" --days 30`,
	RunE: runFetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Linkding connection flags
	fetchCmd.Flags().String("url", "", "Linkding instance URL (required)")
	fetchCmd.Flags().String("token", "", "Linkding API token (required)")
	fetchCmd.Flags().Duration("timeout", 30*time.Second, "HTTP request timeout")

	// Time range flags
	fetchCmd.Flags().Int("days", 7, "Number of days to fetch (from now)")
	fetchCmd.Flags().String("since", "", "Fetch bookmarks added since this date (YYYY-MM-DD)")
	fetchCmd.Flags().String("until", "", "Fetch bookmarks added until this date (YYYY-MM-DD)")

	// Filter flags
	fetchCmd.Flags().String("query", "", "Search query to filter bookmarks")

	// Output flags
	fetchCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	fetchCmd.Flags().String("title", "Bookmarks", "Title for the markdown document")
	fetchCmd.Flags().Bool("no-notes", false, "Exclude notes from output")
	fetchCmd.Flags().Bool("no-tags", false, "Exclude tags from output")
	fetchCmd.Flags().Bool("no-group-by-date", false, "Don't group bookmarks by date")
	fetchCmd.Flags().String("date-format", "2006-01-02", "Date format for grouping (Go time format)")
	fetchCmd.Flags().String("template", "", "Custom template file (default: built-in template)")

	// Bind flags to viper
	_ = viper.BindPFlag("linkding.url", fetchCmd.Flags().Lookup("url"))
	_ = viper.BindPFlag("linkding.token", fetchCmd.Flags().Lookup("token"))
	_ = viper.BindPFlag("linkding.timeout", fetchCmd.Flags().Lookup("timeout"))
	_ = viper.BindPFlag("fetch.days", fetchCmd.Flags().Lookup("days"))
	_ = viper.BindPFlag("fetch.since", fetchCmd.Flags().Lookup("since"))
	_ = viper.BindPFlag("fetch.until", fetchCmd.Flags().Lookup("until"))
	_ = viper.BindPFlag("fetch.query", fetchCmd.Flags().Lookup("query"))
	_ = viper.BindPFlag("fetch.output", fetchCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("fetch.title", fetchCmd.Flags().Lookup("title"))
	_ = viper.BindPFlag("fetch.no_notes", fetchCmd.Flags().Lookup("no-notes"))
	_ = viper.BindPFlag("fetch.no_tags", fetchCmd.Flags().Lookup("no-tags"))
	_ = viper.BindPFlag("fetch.no_group_by_date", fetchCmd.Flags().Lookup("no-group-by-date"))
	_ = viper.BindPFlag("fetch.date_format", fetchCmd.Flags().Lookup("date-format"))
	_ = viper.BindPFlag("fetch.template", fetchCmd.Flags().Lookup("template"))
}

func runFetch(cmd *cobra.Command, args []string) error {
	logger := GetLogger()

	// Get Linkding connection settings
	linkdingURL := viper.GetString("linkding.url")
	linkdingToken := viper.GetString("linkding.token")
	timeout := viper.GetDuration("linkding.timeout")

	if linkdingURL == "" {
		return fmt.Errorf("Linkding URL is required (--url or config file)")
	}
	if linkdingToken == "" {
		return fmt.Errorf("Linkding API token is required (--token or config file)")
	}

	// Create Linkding client
	logger.Infof("Connecting to Linkding at %s", linkdingURL)
	client, err := linkding.NewClient(linkdingURL, linkdingToken, timeout)
	if err != nil {
		return fmt.Errorf("failed to create Linkding client: %w", err)
	}

	// Parse time range
	var addedSince, addedUntil time.Time

	sinceStr := viper.GetString("fetch.since")
	untilStr := viper.GetString("fetch.until")
	days := viper.GetInt("fetch.days")

	if sinceStr != "" {
		addedSince, err = time.Parse("2006-01-02", sinceStr)
		if err != nil {
			return fmt.Errorf("invalid --since date format (use YYYY-MM-DD): %w", err)
		}
	} else if days > 0 {
		addedSince = time.Now().AddDate(0, 0, -days)
	}

	if untilStr != "" {
		addedUntil, err = time.Parse("2006-01-02", untilStr)
		if err != nil {
			return fmt.Errorf("invalid --until date format (use YYYY-MM-DD): %w", err)
		}
	}

	// Fetch bookmarks
	query := viper.GetString("fetch.query")
	logger.Infof("Fetching bookmarks since %s", addedSince.Format("2006-01-02"))
	bookmarks, err := client.FetchAllBookmarks(query, addedSince, addedUntil)
	if err != nil {
		return fmt.Errorf("failed to fetch bookmarks: %w", err)
	}

	logger.Infof("Fetched %d bookmarks", len(bookmarks))

	if len(bookmarks) == 0 {
		logger.Warn("No bookmarks found matching the criteria")
		return nil
	}

	// Create markdown generator
	templatePath := viper.GetString("fetch.template")
	var generator *markdown.Generator
	if templatePath != "" {
		logger.Infof("Using custom template: %s", templatePath)
		generator, err = markdown.NewGeneratorFromFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to create markdown generator from template: %w", err)
		}
	} else {
		generator, err = markdown.NewGenerator()
		if err != nil {
			return fmt.Errorf("failed to create markdown generator: %w", err)
		}
	}

	// Prepare generation options
	opts := markdown.Options{
		Title:          viper.GetString("fetch.title"),
		IncludeNotes:   !viper.GetBool("fetch.no_notes"),
		IncludeTags:    !viper.GetBool("fetch.no_tags"),
		GroupByDate:    !viper.GetBool("fetch.no_group_by_date"),
		DateFormat:     viper.GetString("fetch.date_format"),
	}

	// Determine output destination
	outputPath := viper.GetString("fetch.output")
	var output *os.File
	if outputPath != "" {
		output, err = os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() { _ = output.Close() }()
		logger.Infof("Writing output to %s", outputPath)
	} else {
		output = os.Stdout
	}

	// Generate markdown
	if err := generator.Generate(output, bookmarks, opts); err != nil {
		return fmt.Errorf("failed to generate markdown: %w", err)
	}

	if outputPath != "" {
		logger.Infof("Successfully wrote markdown to %s", outputPath)
	}

	return nil
}
