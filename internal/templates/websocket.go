package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/websocket/server.go
var websocketMain string

type Websocket struct{}

func (t Websocket) Name() string {
	return "websocket"
}

func (t Websocket) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, websocketMain)
}

func (t Websocket) PrintNextSteps() {

}

func init() {
	Register(Websocket{})
}
