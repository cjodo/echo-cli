package templates

import _ "embed"

//go:embed files/embed-resources/server.go
var embedResourcesContent string

func init() {
	Register(NewSingleFileGenerator("embed-resources", embedResourcesContent, "TODO: Implement next steps embed-resources"))
}
