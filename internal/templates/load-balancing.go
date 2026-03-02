package templates

import _ "embed"

//go:embed files/load-balancing/server.go
var loadBalancingContent string

func init() {
	Register(NewSingleFileGenerator("load-balancing", loadBalancingContent, ""))
}
