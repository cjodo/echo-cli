package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func rootRunE(cmd *cobra.Command, _ []string) error {
	return fmt.Errorf("")
}

var rootCmd = &cobra.Command{
	Use:               "echo-cli",
	Long:              "",
	RunE:              rootRunE,
	PersistentPreRunE: preRunE,
	SilenceErrors:     true,
}

func init() {
	rootCmd.AddCommand(
		cookbookCmd,
		docsCmd,
		versionCmd,
		upgradeCmd,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// TODO: Handle properly
	}
}

func preRunE(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "version" {
		return nil
	}

	current := resolveVersion()

	if current == "dev" || len(current) == 0 {
		return nil
	}

	upgradeAvailable, latest := checkForUpgrade()
	if upgradeAvailable {
		fmt.Printf("\n🚀 A new version of echo-cli is available: %s → %s\n", current, latest)
		fmt.Println("Run:")
		fmt.Println("  go install github.com/cjodo/echo-cli@latest")
	}

	return nil
}
