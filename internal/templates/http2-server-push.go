package templates

import (
	"embed"
	_ "embed"
)

//go:embed files/http2-server-push
var http2ServerPushFiles embed.FS

func init() {
	Register(NewMultiFileGenerator("http2-server-push", http2ServerPushFiles, "Implement next steps: http2-server-push"))
}
