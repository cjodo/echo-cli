package templates

import (
	"embed"
	_ "embed"
)

//go:embed files/jsonp
var jsonpFiles embed.FS

func init() {
	Register(NewMultiFileGenerator("jsonp", jsonpFiles, "Implement next steps: jsonp"))
}
