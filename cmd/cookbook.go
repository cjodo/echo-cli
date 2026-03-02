package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/cjodo/echo-cli/internal"
	"github.com/spf13/cobra"
)

type GithubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

const apiCookbookRepo = "https://api.github.com/repos/labstack/echox/contents/cookbook"

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
	url := fmt.Sprintf(apiCookbookRepo+"/%s", recipe)
	outDir := filepath.Join(".", recipe)

	internal.DownloadFromRepo(url, outDir)

	fmt.Printf("Recipe '%s' pulled into %s\n", recipe, outDir)
	fmt.Println("\n\n\n ---Next Steps---\n")
	fmt.Printf("cd %s && go mod tidy\n", recipe)
	return nil
}

func cookbookListRunE(cmd *cobra.Command, args []string) error {
	resp, err := http.Get(apiCookbookRepo)
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
