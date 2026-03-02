package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/reverse-proxy/server.go
var reverseProxyMain string

type ReverseProxy struct{}

func (t ReverseProxy) Name() string {
	return "reverse-proxy"
}

func (t ReverseProxy) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, reverseProxyMain)
}

func (t ReverseProxy) PrintNextSteps() {

}

func init() {
	Register(ReverseProxy{})
}
