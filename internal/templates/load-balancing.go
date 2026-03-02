package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/load-balancing/server.go
var loadBalancingMain string

type LoadBalancing struct{}

func (t LoadBalancing) Name() string {
	return "load-balancing"
}

func (t LoadBalancing) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, loadBalancingMain)
}

func (t LoadBalancing) PrintNextSteps() {

}

func init() {
	Register(LoadBalancing{})
}
