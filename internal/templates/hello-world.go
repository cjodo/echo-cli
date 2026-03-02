package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/hello-world/server.go
var helloWorldMain string

type HelloWorld struct{}

func (t HelloWorld) Name() string {
	return "hello-world"
}

func (t HelloWorld) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, helloWorldMain)
}

func (t HelloWorld) PrintNextSteps() {

}

func init() {
	Register(HelloWorld{})
}
