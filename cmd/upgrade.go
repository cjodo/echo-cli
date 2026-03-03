package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

const latestVersionURL = "https://api.github.com/repos/cjodo/echo-cli/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	Short: "Upgrade echo-cli to latest release",
	RunE: runUpgradeE,
}

func runUpgradeE(cmd *cobra.Command, args []string) error {
	current := resolveVersion()

	// Don't upgrade dev builds
	if current == "dev" {
		fmt.Println("Development build detected. Upgrade skipped.")
		return nil
	}

	fmt.Println("Checking for latest version...")

	upgradeAvailable, latest := checkForUpgrade()
	if !upgradeAvailable {
		fmt.Println("You are already on the latest version.")
		return nil
	}

	fmt.Printf("Upgrading echo-cli from %s → %s...\n\n", current, latest)

	installCmd := exec.Command(
		"go",
		"install",
		"github.com/cjodo/echo-cli@latest",
	)

	installCmd.Stderr = cmd.OutOrStderr() 
	installCmd.Stderr = cmd.ErrOrStderr()

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	fmt.Println("\n✅ Upgrade complete.")
	fmt.Println("Restart your shell if the new version is not immediately available.")

	return nil
}


// checkForUpgrade returns:
//   - false if dev build
//   - false if error
//   - true if a newer version exists
func checkForUpgrade() (bool, string) {
	current := resolveVersion()

	if current == "dev" || strings.HasPrefix(current, "dev") {
		return false, ""
	}

	latest, err := getLatestVersion()
	if err != nil {
		return false, ""
	}

	if isOutdated(current, latest) {
		return true, latest
	}

	return false, ""
}

func getLatestVersion() (string, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(latestVersionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

func isOutdated(current, latest string) bool {
	cv, err := semver.NewVersion(current)
	if err != nil {
		return false
	}

	lv, err := semver.NewVersion(latest)
	if err != nil {
		return false
	}

	return lv.GreaterThan(cv)
}
