package cmd

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// These are overridden by goreleaser via ldflags.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fullVersion())
	},
}

// fullVersion returns the formatted version string shown to users.
func fullVersion() string {
	v := resolveVersion()

	// If this is a release build (commit injected)
	if commit != "" && commit != "none" {
		return fmt.Sprintf("%s (%s, built %s)", v, short(commit), date)
	}

	return v
}

// resolveVersion determines the appropriate version string.
func resolveVersion() string {
	// If GoReleaser injected a version, prefer it.
	if version != "" && version != "dev" {
		return strings.TrimPrefix(version, "v")
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.revision" && len(setting.Value) >= 7 {
			return fmt.Sprintf("dev (%s)", setting.Value[:7])
		}
	}

	return "dev"
}

func short(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}
