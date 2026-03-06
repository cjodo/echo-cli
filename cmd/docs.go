package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	defaultDocsDir        = "/.cache/echo-cli/docs/"
	staticContentRepo      = "https://github.com/cjodo/echo-docs.git"
	docsRefreshCache      bool
	docsPort              string
	verbose               bool
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Serve the echo docs offline + locally",
	RunE:  docsRunE,
}

func init() {
	docsCmd.Flags().BoolVar(&docsRefreshCache, "refresh", false, "Force refresh cache")
	docsCmd.Flags().StringVarP(&docsPort, "port", "p", "8000", "Port to serve docs on")
	docsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

func docsRunE(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	docsPath := home + defaultDocsDir
	_, err = os.Stat(docsPath)

	if docsRefreshCache || os.IsNotExist(err) {
		if docsRefreshCache {
			fmt.Println("Refreshing docs...")
		} else {
			fmt.Println("Cloning docs for offline use...")
		}

		if err := cloneRepo(staticContentRepo, docsPath); err != nil {
			return err
		}

		fmt.Println("Docs ready for offline use")
	}

	fmt.Printf("Serving docs at http://localhost:%s\n", docsPort)
	fs := http.FileServer(http.Dir(docsPath))
	return http.ListenAndServe(":"+docsPort, fs)
}

// cloneRepo runs 'git clone' (or 'git pull' if already exists)
func cloneRepo(repoURL, toPath string) error {
	if _, err := os.Stat(toPath + "/.git"); os.IsNotExist(err) {
		// Clone if repo doesn't exist
		cmd := exec.Command("git", "clone", repoURL, toPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		// Pull if repo exists
		cmd := exec.Command("git", "-C", toPath, "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}
