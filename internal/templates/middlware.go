package templates

import _ "embed"

//go:embed files/middleware/server.go
var middlewareContent string

func init() {
	Register(NewSingleFileGenerator("middleware", middlewareContent, ""))
}
