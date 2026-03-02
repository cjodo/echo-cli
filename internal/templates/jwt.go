package templates

import _ "embed"

//go:embed files/jwt/server.go
var jwtContent string

func init() {
	Register(NewSingleFileGenerator("jwt", jwtContent, ""))
}
