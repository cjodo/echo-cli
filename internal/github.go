package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GithubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

// DownloadFromRepo downloads all files from a GitHub Contents API URL into outDir
// without including the top-level "cookbook" folder.
func DownloadFromRepo(apiURL, outDir string) error {
	parts := strings.Split(apiURL, "/contents/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid GitHub contents API URL: %s", apiURL)
	}
	basePath := parts[1]

	return downloadFromRepoRecursive(apiURL, outDir, basePath)
}

func downloadFromRepoRecursive(apiURL, outDir, basePath string) error {

	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch %s: %s", apiURL, resp.Status)
	}

	var items []GithubContent
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return err
	}

	for _, item := range items {
		switch item.Type {
		case "file":
			if item.DownloadURL == "" {
				continue
			}

			fileResp, err := http.Get(item.DownloadURL)
			if err != nil {
				return err
			}

			if fileResp.StatusCode != http.StatusOK {
				fileResp.Body.Close()
				return fmt.Errorf("failed downloading %s: %s", item.Path, fileResp.Status)
			}

			data, err := io.ReadAll(fileResp.Body)
			fileResp.Body.Close()
			if err != nil {
				return err
			}

			// Strip the "cookbook/..." prefix from item.Path
			relPath, err := filepath.Rel(basePath, item.Path)
			if err != nil {
				return err
			}

			outPath := filepath.Join(outDir, relPath)

			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}

			if err := os.WriteFile(outPath, data, 0644); err != nil {
				return err
			}

			fmt.Println("Downloaded:", outPath)

		case "dir":
			if err := downloadFromRepoRecursive(item.URL, outDir, basePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func ListDirsInRepo(apiUrl string) (*[]GithubContent, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var dirs []GithubContent
	if err := json.Unmarshal(body, &dirs); err != nil {
		return nil, err
	}

	return &dirs, nil
}
