package cmd

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

const latestVersionURL = "https://api.github.com/repos/cjodo/echo-cli/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
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
