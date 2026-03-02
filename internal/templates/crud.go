package templates

import _ "embed"

//go:embed files/crud/server.go
var crudContent string

func init() {
	Register(NewSingleFileGenerator("crud", crudContent, "TODO: Implement next steps crud"))
}
