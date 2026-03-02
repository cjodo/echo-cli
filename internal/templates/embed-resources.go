package templates

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed files/embed-resources/server.go
var embedMain string

type EmbedResources struct{}

func (t EmbedResources) Name() string {
	return "embed-resources"
}

func (t EmbedResources) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return os.WriteFile(mainFile, []byte(embedMain), 0644)
}

func (t EmbedResources) PrintNextSteps() {
	fmt.Println("TODO: Implement next steps embed-resources")
}

func init() {
	Register(EmbedResources{})
}
