package templates

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed files/auto-tls/server.go
var autoTLSMain string

type AutoTLS struct{}

func (t AutoTLS) Name() string {
	return "auto-tls"
}

func (t AutoTLS) Generate(projectPath, modName string) error {
	serverFile := filepath.Join(projectPath, "server.go")
	return os.WriteFile(serverFile, []byte(autoTLSMain), 0644)
}

func(t AutoTLS) PrintNextSteps() {
	fmt.Println("Implement next steps: auto-tls")
}

func init() {
	Register(AutoTLS{})
}
