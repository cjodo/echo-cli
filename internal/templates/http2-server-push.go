package templates

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed files/http2-server-push
var http2ServerPush embed.FS

type Http2ServerPush struct{}

func (t Http2ServerPush) Name() string {
	return "graceful-shutdown"
}

func (t Http2ServerPush) Generate(projectPath, modName string) error {
	return generateFromDir(projectPath, &http2ServerPush, "files/http2-server-push")
}

func(t Http2ServerPush) PrintNextSteps() {
	fmt.Println("Implement next steps: http2-server-push")
}

func init() {
	Register(Http2ServerPush{})
}
