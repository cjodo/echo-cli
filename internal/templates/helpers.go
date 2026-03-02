package templates

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

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
