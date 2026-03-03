package cmd

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current cli version",
	Run:   versionRun,
}

var (
	version string
	unknown = "unknown"
)

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Println(getVersion())
}

func getVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		version = unknown
		return version
	}

	if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		version = strings.TrimPrefix(buildInfo.Main.Version, "v")
		return version
	}

	// Fallback to VCS info if available
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.tag":
			if setting.Value != "" {
				version = strings.TrimPrefix(setting.Value, "v")
				return version
			}
		case "vcs.revision":
			if setting.Value != "" && len(setting.Value) >= 7 {
				version = setting.Value[:7] // short commit hash
				return version
			}
		default:
			continue
		}
	}

	version = unknown
	return version
}

func currentVersion() string {
	return version
}
