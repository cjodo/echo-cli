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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
	SilenceErrors:     true,
}

func init() {
	rootCmd.AddCommand(cookbookCmd, docsCmd, versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {

	}
}
