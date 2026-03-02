package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type GithubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

const repoUrl = "https://api.github.com/repos/labstack/echox/contents/cookbook"

var cookbookCmd =  &cobra.Command{
	Use: "cookbook",
	Short: "Generate a new template from the cookbook",
	Long: "List and pull recipes from the official LabStack EchoX cookbook repo",
}

var cookbookListCmd = &cobra.Command{
	Use: "list",
	Short: "List available cookbook recipies",
	RunE: cookbookListRunE,
}

var cookbookGetCmd = &cobra.Command{
	Use: "get <recipe name>",
	Short: "Pull a recipe from the official cookbook",
	Args: cobra.ExactArgs(1),
	RunE: cookbookGetRunE,
}

func cookbookGetRunE(cmd *cobra.Command, args []string) error {
	recipe := args[0]
	url := fmt.Sprintf(repoUrl+"/%s", recipe)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var files []GithubContent
	if err := json.Unmarshal(body, &files); err != nil {
		return err
	}

	outDir := filepath.Join(".", recipe)
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
	fmt.Printf("Recipe '%s' pulled into %s\n", recipe, outDir)
	return nil
}

func cookbookListRunE(cmd *cobra.Command, args []string) error {
	resp, err := http.Get(repoUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var contents []GithubContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return err
	}

	fmt.Println("Available recipies:")
	for _, c := range contents {
		if c.Type == "dir" {
			fmt.Println(" -", c.Name)
		}
	}
	return nil
}

func init() {
	cookbookCmd.AddCommand(cookbookListCmd, cookbookGetCmd)
}
