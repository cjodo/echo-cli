package templates

import _ "embed"

//go:embed files/websocket/server.go
var websocketContent string

func init() {
	Register(NewSingleFileGenerator("websocket", websocketContent, ""))
}
