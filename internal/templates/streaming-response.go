package templates

import _ "embed"

//go:embed files/streaming-response/server.go
var streamingResponseContent string

func init() {
	Register(NewSingleFileGenerator("streaming-response", streamingResponseContent, ""))
}
