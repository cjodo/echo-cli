package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "echo-cli",
	Long:              "",
	PersistentPreRunE: preRunE,
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
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithColorSchemeFunc(ApplyColorScheme()),
		); err != nil {
		fmt.Println(rootCmd.Help())
	}
}

func preRunE(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "version" {
		return nil
	}

	current := Version()

	if current == "dev" || len(current) == 0 {
		return nil
	}

	upgradeAvailable, latest := checkForUpgrade(current)
	if upgradeAvailable {
		fmt.Printf("\n🚀 A new version of echo-cli is available: %s → %s\n", current, latest)
		fmt.Println("Run:")
		fmt.Println("  go install github.com/cjodo/echo-cli@latest")
	}

	return nil
}

