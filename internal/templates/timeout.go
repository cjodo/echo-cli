package templates

import _ "embed"

//go:embed files/timeout/server.go
var timeoutContent string

func init() {
	Register(NewSingleFileGenerator("timeout", timeoutContent, ""))
}
