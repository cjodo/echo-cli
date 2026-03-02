package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/jwt/server.go
var jwtMain string

type JWT struct{}

func (t JWT) Name() string {
	return "jwt"
}

func (t JWT) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, jwtMain)
}

func (t JWT) PrintNextSteps() {

}

func init() {
	Register(JWT{})
}
