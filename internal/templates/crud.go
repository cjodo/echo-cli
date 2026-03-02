package templates

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed files/crud/server.go
var crudMain string

type Crud struct{}

func (t Crud) Name() string {
	return "crud"
}

func (t Crud) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return os.WriteFile(mainFile, []byte(crudMain), 0644)
}

func (t Crud) PrintNextSteps() {
	fmt.Println("TODO: Implement next steps crud")
}

func init() {
	Register(Crud{})
}
