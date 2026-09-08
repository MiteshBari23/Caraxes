package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitRepo(path string) error {
	gitPath := filepath.Join(path, ".caraxes")

	reinit := false
	if _, err := os.Stat(gitPath); err == nil {
		reinit = true
	}

	dirs := []string{
		filepath.Join(gitPath, "objects"),
		filepath.Join(gitPath, "objects", "info"),
		filepath.Join(gitPath, "objects", "pack"),
		filepath.Join(gitPath, "refs", "heads"),
		filepath.Join(gitPath, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}

	headPath := filepath.Join(gitPath, "HEAD")
	headContent := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(headPath, headContent, 0644); err != nil {
		return fmt.Errorf("writing HEAD: %w", err)
	}

	if reinit {
		fmt.Println("Reinitialized existing Caraxes repository in", gitPath)
	} else {
		fmt.Println("Initialized empty Caraxes repository in", gitPath)
	}

	return nil
}
