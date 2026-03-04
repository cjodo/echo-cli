package internal

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjodo/echo-cli/internal/cache"
)

type FetchOptions struct {
	Verbose bool
	Branch  string
	UseZIP  bool // true for full repo ZIP, false for API directory fetch
}

type RepoSpec struct {
	Host   string
	Owner  string
	Repo   string
	Branch string
	IsGHA  bool
}

type GitContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	SHA         string `json:"sha"`
	URL         string `json:"url"` // <--- add this
}

// ParseRepoSpec parses GitHub/GitLab/Bitbucket repo URLs
func ParseRepoSpec(apiURL string) (*RepoSpec, error) {
	parts := strings.Split(apiURL, "/")
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid repo URL: %s", apiURL)
	}

	spec := &RepoSpec{Branch: "main"}

	switch {
	case strings.Contains(parts[2], "github"):
		spec.Host = "github"
		spec.Owner = parts[4]
		spec.Repo = parts[5]
		spec.IsGHA = strings.Contains(apiURL, "api.github.com")
	case strings.Contains(parts[2], "gitlab"):
		spec.Host = "gitlab"
		spec.Owner = parts[3]
		spec.Repo = parts[4]
	case strings.Contains(parts[2], "bitbucket"):
		spec.Host = "bitbucket"
		spec.Owner = parts[3]
		spec.Repo = parts[4]
	default:
		return nil, fmt.Errorf("unsupported Git host: %s", parts[2])
	}

	return spec, nil
}

func (r *RepoSpec) ZIPURL() string {
	switch r.Host {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", r.Owner, r.Repo, r.Branch)
	case "gitlab":
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/archive/%s/%s-%s.zip", r.Owner, r.Repo, r.Branch, r.Repo, r.Branch)
	case "bitbucket":
		return fmt.Sprintf("https://bitbucket.org/%s/%s/get/%s.zip", r.Owner, r.Repo, r.Branch)
	default:
		return ""
	}
}

func (r *RepoSpec) CacheKey() string {
	return fmt.Sprintf("%s:%s:%s/%s", r.Host, r.Branch, r.Owner, r.Repo)
}

// DownloadFromRepoWithCache downloads a repo as ZIP or recursively via API
func DownloadFromRepoWithCache(apiURL, outDir string, c *cache.Cache, opts FetchOptions) error {
	if opts.UseZIP {
		return downloadRepoZIP(apiURL, outDir, c, opts)
	}
	return downloadFromRepoAPI(apiURL, outDir, c, opts)
}

///////////////////////////
// ZIP-based repo download
///////////////////////////
func downloadRepoZIP(apiURL, outDir string, c *cache.Cache, opts FetchOptions) error {
	spec, err := ParseRepoSpec(apiURL)
	if err != nil {
		return err
	}
	if opts.Branch != "" {
		spec.Branch = opts.Branch
	}

	cacheKey := spec.CacheKey()

	if c != nil {
		if cachedPath, _ := c.GetFile(cacheKey); cachedPath != "" {
			if opts.Verbose {
				fmt.Println("Using cached ZIP for repo")
			}
			return unzip(cachedPath, outDir, opts.Verbose)
		}
	}

	zipURL := spec.ZIPURL()
	if opts.Verbose {
		fmt.Println("Downloading ZIP from:", zipURL)
	}

	resp, err := http.Get(zipURL)
	if err != nil {
		return fmt.Errorf("failed to download ZIP: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download ZIP: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "repo-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write ZIP: %w", err)
	}
	tmpFile.Seek(0, io.SeekStart)

	if c != nil {
		if _, err := c.SetFile(cacheKey, []byte(tmpFile.Name())); err != nil && opts.Verbose {
			fmt.Println("Warning: failed to cache ZIP:", err)
		}
	}

	return unzip(tmpFile.Name(), outDir, opts.Verbose)
}

func unzip(zipPath, outDir string, verbose bool) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) < 2 {
			continue
		}
		relPath := filepath.Join(outDir, parts[1])

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(relPath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(relPath), 0755); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(relPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			dstFile.Close()
			srcFile.Close()
			return err
		}
		dstFile.Close()
		srcFile.Close()

		if verbose {
			fmt.Println("Extracted:", relPath)
		}
	}

	return nil
}

///////////////////////////
// API recursive download
///////////////////////////
func downloadFromRepoAPI(apiURL, outDir string, c *cache.Cache, opts FetchOptions) error {
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

	var items []GitContent
	if err := json.Unmarshal(body, &items); err != nil {
		return err
	}

	for _, item := range items {
		switch item.Type {
		case "file":
			if item.DownloadURL == "" {
				return fmt.Errorf("missing download URL for file %s", item.Path)
			}
			data, err := downloadFile(item, c)
			if err != nil {
				return err
			}
			outPath := filepath.Join(outDir, item.Path)
			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(outPath, data, 0644); err != nil {
				return err
			}
			if opts.Verbose {
				fmt.Println("Downloaded:", outPath)
			}
		case "dir":
			if err := downloadFromRepoAPI(item.URL, outDir, c, opts); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadFile(item GitContent, c *cache.Cache) ([]byte, error) {
	if c != nil && item.DownloadURL != "" {
		if cachedPath, err := c.GetFile(item.DownloadURL); err == nil {
			if data, err := os.ReadFile(cachedPath); err == nil {
				return data, nil
			}
		}
	}

	resp, err := http.Get(item.DownloadURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed downloading %s: %s", item.Path, resp.Status)
	}
	data, _ := io.ReadAll(resp.Body)
	if c != nil && item.DownloadURL != "" {
		c.SetFile(item.DownloadURL, data)
	}
	return data, nil
}


///////////////////////////
// List directories in repo
///////////////////////////
func ListDirsInRepoWithCache(apiURL string, c *cache.Cache) (*[]GitContent, error) {
	cacheKey := "list:" + apiURL
	if c != nil {
		if data, err := c.GetAPI(cacheKey); err == nil {
			var contents []GitContent
			if err := json.Unmarshal(data, &contents); err == nil {
				return &contents, nil
			}
		}
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch directory listing: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch directory listing: %s", resp.Status)
	}

	var contents []GitContent
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if c != nil {
		if data, err := json.Marshal(contents); err == nil {
			_ = c.SetAPI(cacheKey, data)
		}
	}

	return &contents, nil
}
