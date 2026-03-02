package templates

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type SingleFileGenerator struct {
	name           string
	content        string
	printNextSteps string
}

func (s SingleFileGenerator) Name() string {
	return s.name
}

func (s SingleFileGenerator) Generate(projectPath, modName string) error {
	mainFile := filepath.Join(projectPath, "server.go")
	return os.WriteFile(mainFile, []byte(s.content), 0644)
}

func (s SingleFileGenerator) PrintNextSteps() {
	if s.printNextSteps != "" {
		fmt.Println(s.printNextSteps)
	}
}

func NewSingleFileGenerator(name, content, printNextSteps string) SingleFileGenerator {
	return SingleFileGenerator{
		name:           name,
		content:        content,
		printNextSteps: printNextSteps,
	}
}

type MultiFileGenerator struct {
	name           string
	files          embed.FS
	printNextSteps string
}

func (m MultiFileGenerator) Name() string {
	return m.name
}

func (m MultiFileGenerator) Generate(projectPath, modName string) error {
	return generateFromDir(projectPath, &m.files, fmt.Sprintf("files/%s", m.name))
}

func (m MultiFileGenerator) PrintNextSteps() {
	if m.printNextSteps != "" {
		fmt.Println(m.printNextSteps)
	}
}

func NewMultiFileGenerator(name string, files embed.FS, printNextSteps string) MultiFileGenerator {
	return MultiFileGenerator{
		name:           name,
		files:          files,
		printNextSteps: printNextSteps,
	}
}

func generateFromDir(projectPathRoot string, efs *embed.FS, root string) error {
	return fs.WalkDir(*efs, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root folder itself
		if path == root {
			return nil
		}

		// Compute relative path from the root folder
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// Destination path: flatten so only the files under root appear in projectPathRoot
		targetPath := filepath.Join(projectPathRoot, relPath)

		if d.IsDir() {
			// Create any nested directories if needed (e.g., subfolders in the template)
			return os.MkdirAll(targetPath, 0755)
		}

		// Ensure parent directories exist
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Open source embedded file
		srcFile, err := efs.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create destination file
		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// Copy contents
		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

func generateSingle(name, content string) error {
	return os.WriteFile(name, []byte(content), 0644)
}
