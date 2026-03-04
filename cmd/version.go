package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	// These are overridden by goreleaser via ldflags.
	// For testing purposes
	skipGitVersion bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(Version())
	},
}

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	return "dev"
}
