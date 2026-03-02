package templates

import (
	"embed"
	_ "embed"
)

//go:embed files/file-download
var fileDownloadFiles embed.FS

func init() {
	Register(NewMultiFileGenerator("file-download", fileDownloadFiles, "TODO: Implement next steps file download"))
}
