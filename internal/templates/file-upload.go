package templates

import (
	"embed"
	_ "embed"
)

//go:embed files/file-upload
var fileUploadFiles embed.FS

func init() {
	Register(NewMultiFileGenerator("file-upload", fileUploadFiles, "TODO: Implement next steps file upload"))
}
