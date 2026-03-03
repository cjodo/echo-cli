package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjodo/echo-cli/internal/cache"
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
	return DownloadFromRepoWithCache(apiURL, outDir, nil)
}

func DownloadFromRepoWithCache(apiURL, outDir string, c *cache.Cache) error {
	parts := strings.Split(apiURL, "/contents/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid GitHub contents API URL: %s", apiURL)
	}
	basePath := parts[1]

	return downloadFromRepoRecursive(apiURL, outDir, basePath, c)
}

func downloadFromRepoRecursive(apiURL, outDir, basePath string, c *cache.Cache) error {
	var body []byte
	var err error

	if c != nil {
		body, err = c.FetchWithCache(apiURL)
	} else {
		resp, err := http.Get(apiURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to fetch %s: %s", apiURL, resp.Status)
		}
		body, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		return err
	}

	var items []GithubContent
	if err := json.Unmarshal(body, &items); err != nil {
		return err
	}

	for _, item := range items {
		switch item.Type {
		case "file":
			if item.DownloadURL == "" {
				continue
			}

			data, err := func() ([]byte, error) {
				if c != nil {
					cachedPath, err := c.GetFile(item.DownloadURL)
					if err == nil {
						data, err := os.ReadFile(cachedPath)
						if err == nil {
							return data, nil
						}
					}
				}

				fileResp, err := http.Get(item.DownloadURL)
				if err != nil {
					return nil, err
				}
				defer fileResp.Body.Close()

				if fileResp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("failed downloading %s: %s", item.Path, fileResp.Status)
				}

				return io.ReadAll(fileResp.Body)
			}()

			if err != nil {
				return err
			}

			if c != nil {
				_, _ = c.SetFile(item.DownloadURL, data)
			}
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
			if err := downloadFromRepoRecursive(item.URL, outDir, basePath, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func ListDirsInRepo(apiUrl string) (*[]GithubContent, error) {
	return ListDirsInRepoWithCache(apiUrl, nil)
}

func ListDirsInRepoWithCache(apiURL string, c *cache.Cache) (*[]GithubContent, error) {
	var body []byte
	var err error

	if c != nil {
		body, err = c.FetchWithCache(apiURL)
	} else {
		resp, err := http.Get(apiURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		return nil, err
	}

	var dirs []GithubContent
	if err := json.Unmarshal(body, &dirs); err != nil {
		return nil, err
	}

	return &dirs, nil
}
