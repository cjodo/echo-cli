package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cjodo/echo-cli/internal"
	"github.com/cjodo/echo-cli/internal/cache"
	"github.com/spf13/cobra"
)

var (
	defaultDocsDir        = "/.echo-cli/docs/"
	staticContentRepoBase = "https://api.github.com/repos/cjodo/echo-docs/contents"
	docsCache             *cache.Cache
	docsRefreshCache      bool
	docsPort              string
	verbose               bool
)

func init() {
	var err error
	docsCache, err = cache.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: docs cache init failed: %v\n", err)
	}
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Serve the echo docs offline + locally",
	RunE:  docsRunE,
}

func init() {
	docsCmd.Flags().BoolVar(&docsRefreshCache, "refresh", false, "Force refresh cache")
	docsCmd.Flags().StringVarP(&docsPort, "port", "p", "8080", "Port to serve docs on")
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
			fmt.Println("Refreshing docs cache...")
		} else {
			fmt.Println("Downloading docs for offline use...")
		}

		if err := downloadDocs(docsPath); err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("docs cache ready for offline use")
	}

	fmt.Printf("Serving docs at http://localhost:%s\n", docsPort)

	fs := http.FileServer(http.Dir(docsPath))
	return http.ListenAndServe(":"+docsPort, fs)
}

func downloadDocs(toPath string) error {
	c := docsCache
	if docsRefreshCache {
		c = nil
	}

	opts := internal.FetchOptions{
		Verbose: verbose,
		UseZIP: true,
	}

	if err := internal.DownloadFromRepoWithCache(staticContentRepoBase, toPath, c, opts); err != nil {
		return err
	}

	return nil
}
