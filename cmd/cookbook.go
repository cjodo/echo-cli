package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	Use: "cookbook <template>",
	Aliases: []string{"c"},
	Short: "Generate a new template from the cookbook",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return nil },
}

var cookbookListCmd = &cobra.Command{
	Use: "list",
	Short: "List available cookbook recipies",
	RunE: cookbookListRunE,
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
	cookbookCmd.AddCommand(cookbookListCmd)
}
