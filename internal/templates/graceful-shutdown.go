package templates

import _ "embed"

//go:embed files/graceful-shutdown/server.go
var gracefulShutdownContent string

func init() {
	Register(NewSingleFileGenerator("graceful-shutdown", gracefulShutdownContent, "Implement next steps: graceful-shutdown"))
}
