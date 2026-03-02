package templates

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed files/jsonp
var jsonPDir embed.FS

type JsonP struct{}

func (t JsonP) Name() string {
	return "jsonp"
}

func (t JsonP) Generate(projectPath, modName string) error {
	return generateFromDir(projectPath, &jsonPDir, "files/jsonp")
}

func (t JsonP) PrintNextSteps() {
	fmt.Println("Implement next steps: jsonp")
}

func init() {
	Register(JsonP{})
}
