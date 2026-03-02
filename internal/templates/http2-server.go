package templates

import _ "embed"

//go:embed files/http2-server/server.go
var http2ServerContent string

func init() {
	Register(NewSingleFileGenerator("http2-server", http2ServerContent, ""))
}
