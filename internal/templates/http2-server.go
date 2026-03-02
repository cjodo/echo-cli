package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/http2-server/server.go
var http2ServerMain string

type Http2Server struct{}

func (t Http2Server) Name() string {
	return "http2-server"
}

func (t Http2Server) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, http2ServerMain)
}

func (t Http2Server) PrintNextSteps() {

}

func init() {
	Register(Http2Server{})
}
