package templates

import _ "embed"

//go:embed files/reverse-proxy/server.go
var reverseProxyContent string

func init() {
	Register(NewSingleFileGenerator("reverse-proxy", reverseProxyContent, ""))
}
