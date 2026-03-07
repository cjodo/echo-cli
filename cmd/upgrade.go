package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// githubRelease represents a GitHub release
type githubRelease struct {
	TagName string `json:"tag_name"`
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade echo-cli to the latest release or Go module version",
	RunE:  runUpgrade,
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	current := Version()
	// if current == "dev" {
	// 	fmt.Println("Development build detected. Upgrade skipped.")
	// 	return nil
	// }

	fmt.Println("Checking for latest version...")

	// Try to get latest GitHub release
	latest, err := getLatestReleaseTag()
	if err != nil {
		fmt.Println("⚠ Could not fetch GitHub release:", err)
		fmt.Println("Falling back to installing latest Go module version...")
		latest = "latest"
	}

	upgradeMsg := latest
	if latest != "latest" {
		upgradeMsg = fmt.Sprintf("v%s", latest)
	}

	fmt.Printf("Upgrading echo-cli from %s → %s...\n\n", current, upgradeMsg)

	installCmd := exec.Command("go", "install", "-mod=mod", fmt.Sprintf("github.com/cjodo/echo-cli@v%s", latest))
	installCmd.Stdout = cmd.OutOrStdout()
	installCmd.Stderr = cmd.ErrOrStderr()
	installCmd.Env = append(os.Environ(), "GO111MODULE=on")

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	// Check if GOBIN/bin is in PATH
	binPath := filepath.Dir(os.Args[0])
	if !strings.Contains(os.Getenv("PATH"), binPath) {
		fmt.Printf("\n⚠ WARNING: %s may not be in your PATH. Add it to run the new version.\n", binPath)
	}

	fmt.Println("\n✅ Upgrade complete.")
	fmt.Println("Restart your shell if the new version is not immediately available.")

	return nil
}

// getLatestReleaseTag fetches the latest GitHub release tag
func getLatestReleaseTag() (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(latestReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	tag := strings.TrimPrefix(release.TagName, "v")
	if tag == "" {
		return "", fmt.Errorf("empty release tag")
	}

	return tag, nil
}

func checkForUpgrade(current string) (bool, string) {
	return false, ""
}
