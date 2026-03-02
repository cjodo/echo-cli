package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type GithubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

func DownloadFromRepo(apiUrl, outDir string) error {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var files []GithubContent
	if err := json.Unmarshal(body, &files); err != nil {
		return err
	}

	os.MkdirAll(outDir, 0755)

	for _, f := range files {
		if f.Type == "file" && f.DownloadURL != "" {
			fileResp, err := http.Get(f.DownloadURL)
			if err != nil {
				return err
			}
			defer fileResp.Body.Close()

			data, _ := io.ReadAll(fileResp.Body)
			outPath := filepath.Join(outDir, f.Name)
			os.WriteFile(outPath, data, 0644)
			fmt.Println("Downloaded", f.Name)
		}
	}
	return nil

}
	
func ListAllRepoContents(apiUrl string) (*[]GithubContent, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var contents []GithubContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return nil, err
	}

	return &contents, nil
}
