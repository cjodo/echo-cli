package templates

import (
	_ "embed"
	"path/filepath"
)

//go:embed files/streaming-response/server.go
var streamingResponseMain string

type StreamingResponse struct{}

func (t StreamingResponse) Name() string {
	return "streaming-response"
}

func (t StreamingResponse) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return generateSingle(mainFile, streamingResponseMain)
}

func (t StreamingResponse) PrintNextSteps() {

}

func init() {
	Register(StreamingResponse{})
}
