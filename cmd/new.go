package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)


var (
	// sets Unix permissions to rwxr-x---. The owner has full read, write, and execute permissions
	modeRWE = os.FileMode(0750)

	newCmd = &cobra.Command{
		Use: "new <project-name>",
		Aliases: []string{"n"},
		Short: "Generate a new echo project",
		Args: cobra.MinimumNArgs(1),
		RunE: newRunE,
	}
)

func newRunE(cmd *cobra.Command, args [] string) (err error) {

	projectName := args[0]
	modName := projectName
	if len(args) > 1 {
		modName = args[1]
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("err Getwd(): %v", err)
	}

	projectPath := fmt.Sprintf("%s%c%s", wd, os.PathSeparator, projectName)
	if err := createProjectFromPath(projectPath); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rmErr := os.RemoveAll(projectPath); rmErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "failed to remove project dir: %v", rmErr)
			}
		}
	}()

	defer func() {
		// success message
		if err == nil {
			cmd.Printf(nextSteps, projectName)
		}
	}()

	return create(projectPath, modName)
}

func createProjectFromPath(path string) error {
	if err := os.Mkdir(path, modeRWE); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("change dir: %w", err)
	}

	return nil
}

func create(projectPath, modName string) error {
	if err := createFile(fmt.Sprintf("%s%cserver.go", projectPath, os.PathSeparator), template); err != nil {
		return err
	}

	return runCmd(exec.Command("go", "mod", "init", modName))
}

func createFile(filePath, content string) error {
	f, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("create %s: %w", filePath, err)
	}

	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "close file: %v", cerr)
		}
	}()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}

	return nil
}

func runCmd(cmd *exec.Cmd) (err error) {
	var (
		stderr io.ReadCloser
		stdout io.ReadCloser
	)

	if stderr, err = cmd.StderrPipe(); err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}
	go func() {
		if _, cErr := io.Copy(os.Stderr, stderr); cErr != nil {
			fmt.Fprintf(os.Stderr, "copy stderr: %v", cErr)
		}
	}()

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	go func() {
		if _, cErr := io.Copy(os.Stdout, stdout); cErr != nil {
			fmt.Fprintf(os.Stderr, "copy stdout: %v", cErr)
		}
	}()

	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("failed to run %s: %w", cmd.String(), err)
	}

	return err
}

var ( 
	template = `package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}`

	nextSteps = `To get started with your project, run:
	cd %s
	go mod init
	go run server.go
	Browse to http://localhost:1323 and you should see Hello, World! on the page.
	Visit echo's docs here: https://echo.labstack.com/`
)
