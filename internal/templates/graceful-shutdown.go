package templates

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed files/graceful-shutdown/server.go
var gracefulShutdownMain string

type GracefulShutdown struct{}

func (t GracefulShutdown) Name() string {
	return "graceful-shutdown"
}

func (t GracefulShutdown) Generate(projectPath, modName string) error {
	serverFile := filepath.Join(projectPath, "server.go")
	return os.WriteFile(serverFile, []byte(gracefulShutdownMain), 0644)
}

func (t GracefulShutdown) PrintNextSteps() {
	fmt.Println("Implement next steps: graceful-shutdown")
}

func init() {
	Register(GracefulShutdown{})
}
