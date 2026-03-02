package templates

import _ "embed"

//go:embed files/hello-world/server.go
var helloWorldContent string

func init() {
	Register(NewSingleFileGenerator("hello-world", helloWorldContent, ""))
}
