package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/subdomain/server.go
var subdomainMain string

type Subdomain struct{}

func (t Subdomain) Name() string {
	return "subdomain"
}

func (t Subdomain) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, subdomainMain)
}

func (t Subdomain) PrintNextSteps() {

}

func init() {
	Register(Subdomain{})
}
