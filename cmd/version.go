package cmd

import (
	"fmt"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// These are overridden by goreleaser via ldflags.
	Version string

	// For testing purposes
	skipGitVersion bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(fullVersion())
	},
}

// fullVersion returns the formatted version string shown to users.
func fullVersion() string {
	v := resolveVersion()
	return v
}

// resolveVersion determines the appropriate version string.
func resolveVersion() string {
	// If GoReleaser injected a version, prefer it.
	if Version != "" && Version != "dev" {
		return strings.TrimPrefix(Version, "v")
	}

	// Try to get version from git tags for dev builds (skip in tests)
	if !skipGitVersion {
		if version := getGitVersion(); version != "" {
			return version
		}
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

// getGitVersion tries to get the most recent git tag for dev builds.
func getGitVersion() string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return ""
	}
	version := strings.TrimSpace(string(out))
	if version == "" || strings.HasPrefix(version, "v") {
		return strings.TrimPrefix(version, "v")
	}
	return version
}

func short(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}
