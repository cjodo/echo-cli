package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cjodo/echo-cli/internal"
	"github.com/cjodo/echo-cli/internal/cache"
	"github.com/spf13/cobra"
)

var (
	cookbookCache   *cache.Cache
	refreshCache    bool
	cookbookVerbose bool
)

func init() {
	var err error
	cookbookCache, err = cache.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cache init failed: %v\n", err)
	}

	cookbookGetCmd.Flags().BoolVar(&refreshCache, "refresh", false, "Force refresh cache")
	cookbookGetCmd.Flags().BoolVarP(&cookbookVerbose, "verbose", "v", false, "Enable verbose output")
	cookbookListCmd.Flags().BoolVar(&refreshCache, "refresh", false, "Force refresh cache")
	cookbookListCmd.Flags().BoolVarP(&cookbookVerbose, "verbose", "v", false, "Enable verbose output")
}

type GithubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

const apiCookbookRepo = "https://api.github.com/repos/labstack/echox/contents/cookbook"

var cookbookCmd = &cobra.Command{
	Use:   "cookbook",
	Short: "Generate a new template from the cookbook",
	Long:  "List and pull recipes from the official LabStack EchoX cookbook repo",
}

var cookbookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available cookbook recipies",
	RunE:  cookbookListRunE,
}

var cookbookGetCmd = &cobra.Command{
	Use:   "get <recipe name>",
	Short: "Pull a recipe from the official cookbook",
	Args:  cobra.ExactArgs(1),
	RunE:  cookbookGetRunE,
}

func cookbookGetRunE(cmd *cobra.Command, args []string) error {
	recipe := args[0]
	url := fmt.Sprintf(apiCookbookRepo+"/%s", recipe)
	outDir := filepath.Join(".", recipe)

	c := cookbookCache
	if refreshCache {
		c = nil
	}

	if err := internal.DownloadFromRepoWithCache(url, outDir, c, internal.Options{Verbose: cookbookVerbose}); err != nil {
		return err
	}

	if err := os.Chdir(outDir); err != nil {
		return err
	}

	if err := exec.Command("go", "mod", "init", recipe).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod init failed: %v\n", err)
	}

	if err := os.Chdir(".."); err != nil {
		return err
	}

	fmt.Printf("Recipe '%s' pulled into %s\n", recipe, outDir)
	fmt.Println("---Next Steps---")
	fmt.Printf("cd %s && go mod tidy\n\n", recipe)
	return nil
}

func cookbookListRunE(cmd *cobra.Command, args []string) error {
	c := cookbookCache
	if refreshCache {
		c = nil
	}

	contents, err := internal.ListDirsInRepoWithCache(apiCookbookRepo, c)
	if err != nil {
		return err
	}

	fmt.Println("Available recipies:")
	for _, c := range *contents {
		if c.Type == "dir" {
			fmt.Println(" -", c.Name)
		}
	}
	return nil
}

func init() {
	cookbookCmd.AddCommand(cookbookListCmd, cookbookGetCmd)
}
