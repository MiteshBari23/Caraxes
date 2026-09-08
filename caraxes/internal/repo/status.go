package repo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

var headerStyle = lipgloss.NewStyle().Bold(true)
var hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true) //245 grey
var modifiedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) //196 red
var untrackedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // 214 orange

func StatusRepo() error {
	indexPath := filepath.Join(".caraxes", "index")

	entries, err := readIndex(indexPath)
	if err != nil {
		return err
	}

	var untracked []string
	var modified []string

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".caraxes" {
				return filepath.SkipDir
			}
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hashString, _, err := hashBlob(content)
		if err != nil {
			return err
		}

		relativePath := normalizePath(path)

		storedHash, exists := entries[relativePath]
		if !exists {
			untracked = append(untracked, relativePath)
		} else if storedHash != hashString {
			modified = append(modified, relativePath)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(untracked) == 0 && len(modified) == 0 {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	if len(modified) > 0 {
		fmt.Println(headerStyle.Render("Changes not staged for commit:"))
		fmt.Println(hintStyle.Render("  (use \"caraxes add <file>...\" to update what will be committed)"))
		for _, path := range modified {
			fmt.Println(modifiedStyle.Render("        modified:   " + path))
		}
		fmt.Println()
	}

	if len(untracked) > 0 {
		fmt.Println(headerStyle.Render("Untracked files:"))
		fmt.Println(hintStyle.Render("  (use \"caraxes add <file>...\" to include in what will be committed)"))
		for _, path := range untracked {
			fmt.Println(untrackedStyle.Render("        " + path))
		}
	}

	return nil
}
