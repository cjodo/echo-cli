package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const latestReleaseURL = "https://api.github.com/repos/cjodo/echo-cli/releases/latest"

// versionCmd prints the current CLI version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		// Always show the full version
		fmt.Println(FullVersion())

		// If this is a dev build, optionally show the latest release
		if Version() == "dev" {
			fmt.Println("Latest release:", DevWithLatest())
		}
	},
}
// Version returns the current CLI version.
// Tagged builds return the tag dev builds return "dev".
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}

	return "dev"
}

// LatestRelease returns the latest GitHub release tag (0.1.7)
func LatestRelease() (string, error) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(latestReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	return strings.TrimPrefix(r.TagName, "v"), nil
}

// FullVersion returns the string displayed to users.
func FullVersion() string {
	v := Version()

	if v == "dev" {
		// Include commit hash if available
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, s := range buildInfo.Settings {
				if s.Key == "vcs.revision" && len(s.Value) >= 7 {
					return fmt.Sprintf("dev (%s)", s.Value[:7])
				}
			}
		}
		return "dev"
	}

	return v
}

// Optional helper to show the latest version in a dev build
func DevWithLatest() string {
	v := Version()
	if v != "dev" {
		return v
	}

	latest, err := LatestRelease()
	if err != nil {
		return v
	}

	return fmt.Sprintf("dev (latest %s)", latest)
}
