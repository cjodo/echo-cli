package templates

import _ "embed"

//go:embed files/cors/server.go
var corsContent string

func init() {
	Register(NewSingleFileGenerator("cors", corsContent, "TODO: Implement next steps cors"))
}
