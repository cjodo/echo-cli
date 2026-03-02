package templates

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed files/file-upload
var fileUploadDir embed.FS

type FileUpload struct{}

func (t FileUpload) Name() string {
	return "file-upload"
}

func (t FileUpload) Generate(projectPath, modName string) error {
	return generateFromDir(projectPath, &fileUploadDir, "files/file-upload")
}

func (t FileUpload) PrintNextSteps() {
	fmt.Println("TODO: Implement next steps file upload")
}

func init() {
	Register(FileUpload{})
}
