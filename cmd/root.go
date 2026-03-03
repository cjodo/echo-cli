package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func rootRunE(cmd *cobra.Command, _ []string) error {
	return fmt.Errorf("help: %w", cmd.Help())
}

var rootCmd = &cobra.Command{
	Use:               "echo-cli",
	Long:              "",
	RunE:              rootRunE,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
	SilenceErrors:     true,
}

func init() {
	rootCmd.AddCommand(cookbookCmd, docsCmd, versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error executing echo cli: ", err)
	}
}
