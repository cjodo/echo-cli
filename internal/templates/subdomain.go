package templates

import _ "embed"

//go:embed files/subdomain/server.go
var subdomainContent string

func init() {
	Register(NewSingleFileGenerator("subdomain", subdomainContent, ""))
}
