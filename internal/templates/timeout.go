package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/timeout/server.go
var timeoutMain string

type Timeout struct{}

func (t Timeout) Name() string {
	return "timeout"
}

func (t Timeout) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, timeoutMain)
}

func (t Timeout) PrintNextSteps() {

}

func init() {
	Register(Timeout{})
}
