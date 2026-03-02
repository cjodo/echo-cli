package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/middleware/server.go
var middlewareMain string

type Middleware struct{}

func (t Middleware) Name() string {
	return "middleware"
}

func (t Middleware) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, middlewareMain)
}

func (t Middleware) PrintNextSteps() {

}

func init() {
	Register(Middleware{})
}
