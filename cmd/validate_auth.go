package cmd

import (
	"fmt"
	"time"

	"github.com/lmorchard/linkding-to-markdown/internal/linkding"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateAuthCmd answers "do my credentials work?" with a single-line
// stdout and a clean exit code. Suitable for orchestrators or scripts
// gating behavior on auth health.
var validateAuthCmd = &cobra.Command{
	Use:   "validate-auth",
	Short: "Check whether the configured Linkding URL + token are accepted",
	Long: `Run a minimal authenticated request against the configured Linkding
instance and exit 0 if the credentials are accepted, non-zero otherwise.

Reads ` + "`url`" + ` and ` + "`token`" + ` from the usual sources
(flags / LINKDING_URL + LINKDING_TOKEN env vars / config file).`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := viper.GetString("url")
		token := viper.GetString("token")
		timeout := viper.GetDuration("timeout")
		if timeout == 0 {
			timeout = 30 * time.Second
		}

		if url == "" {
			return fmt.Errorf("linkding URL not configured (set LINKDING_URL or `url:`)")
		}
		if token == "" {
			return fmt.Errorf("linkding API token not configured (set LINKDING_TOKEN or `token:`)")
		}

		client, err := linkding.NewClient(url, token, timeout)
		if err != nil {
			return err
		}
		if err := client.ValidateAuth(); err != nil {
			return err
		}
		fmt.Printf("validate-auth: ok (%s)\n", url)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateAuthCmd)
}
