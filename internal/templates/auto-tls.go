package templates

import _ "embed"

//go:embed files/auto-tls/server.go
var autoTLSContent string

func init() {
	Register(NewSingleFileGenerator("auto-tls", autoTLSContent, "Implement next steps: auto-tls"))
}
