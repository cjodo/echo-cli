package templates

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed files/file-download
var fileDownloadDir embed.FS

type FileDownload struct{}

func (t FileDownload) Name() string {
	return "file-download"
}

func (t FileDownload) Generate(projectPath, modName string) error {
	return generateFromDir(projectPath, &fileDownloadDir, "files/file-download")
}

func (t FileDownload) PrintNextSteps() {
	fmt.Println("TODO: Implement next steps file download")
}

func init() {
	Register(FileDownload{})
}
