package templates

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed files/cors/server.go
var corsMain string

type Cors struct{}

func (t Cors) Name() string {
	return "cors"
}

func (t Cors) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return os.WriteFile(mainFile, []byte(corsMain), 0644)
}

func (t Cors) PrintNextSteps() {
	fmt.Println("TODO: Implement next steps cors")
}

func init() {
	Register(Cors{})
}
