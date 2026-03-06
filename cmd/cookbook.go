package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	cookbookVerbose bool
)

func init() {
	cookbookGetCmd.Flags().BoolVarP(&cookbookVerbose, "verbose", "v", false, "Enable verbose output")
	cookbookListCmd.Flags().BoolVarP(&cookbookVerbose, "verbose", "v", false, "Enable verbose output")

	cookbookCmd.AddCommand(cookbookGetCmd, cookbookListCmd)
}

var cookbookCmd = &cobra.Command{
	Use:   "cookbook",
	Short: "Generate a new template from the cookbook",
	Long:  "List and pull recipes from the official LabStack EchoX cookbook repo",
}

var cookbookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available cookbook recipies",
	RunE:  func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available recipes: ")

		for _, recipe := range recipes {
			fmt.Println(recipe)
		}

		fmt.Println("For more info visit: https://echo.labstack.com/docs/category/cookbook")
		return nil
	},
}

var cookbookGetCmd = &cobra.Command{
	Use:   "get <recipe name>",
	Short: "Pull a recipe from the official cookbook",
	Args:  cobra.ExactArgs(1),
	RunE:  cookbookGetRunE,
}

func cookbookGetRunE(cmd *cobra.Command, args []string) error {
	start := time.Now()
	recipe := args[0]

	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	outDir := filepath.Join(".", recipe)

	if cookbookVerbose {
		fmt.Println("Downloading recipe:", recipe)
	}

	repoUrl := fmt.Sprintf("git@github.com:recipes-echo/%s.git", recipe)
	if err := exec.Command(git, "clone", repoUrl).Run(); err != nil {
		return err
	}

	// Initialize Go module in the recipe folder
	if err := os.Chdir(outDir); err != nil {
		return err
	}
	if err := exec.Command("go", "mod", "init", recipe).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod init failed: %v\n", err)
	}
	if err := os.Chdir(".."); err != nil {
		return err
	}

	elapsed := time.Since(start)

	fmt.Printf("✅ Recipe '%s' pulled into %s\n", recipe, outDir)
	fmt.Println("---Next Steps---")
	fmt.Printf("cd %s && go mod tidy\n\n", recipe)
	fmt.Printf("took: %s\n", elapsed.Round(time.Millisecond))
	return nil
}

var recipes = []string{
    "websocket",
    "timeout",
    "subdomain",
    "streaming-response",
    "sse",
    "reverse-proxy",
    "prometheus",
    "middleware",
    "load-balancing",
    "jwt",
    "jsonp",
    "http2-server-push",
    "http2",
    "hello-world",
    "graceful-shutdown",
    "file-upload",
    "file-download",
    "embed",
    "csrf",
    "crud",
    "cors",
    "casbin",
    "auto-tls",
}
