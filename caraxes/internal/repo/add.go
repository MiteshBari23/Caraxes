package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func AddRepo(path string) error {
	objectsDir := filepath.Join(".caraxes", "objects")
	indexPath := filepath.Join(".caraxes", "index")

	entries, err := readIndex(indexPath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("File does not exists")
		}
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	if fileInfo.IsDir() {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == ".caraxes" {
					return filepath.SkipDir
				}
				return nil
			}
			content, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			hashString, data, err := hashBlob(content)
			if err != nil {
				return err
			}
			if err := writeObject(objectsDir, hashString, data); err != nil {
				return err
			}
			entries[normalizePath(p)] = hashString

			fmt.Println(hashString)
			return nil
		})
		if err != nil {
			return err
		}

		return writeIndex(indexPath, entries)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	hashString, data, err := hashBlob(content)
	if err != nil {
		return err
	}
	if err := writeObject(objectsDir, hashString, data); err != nil {
		return err
	}

	entries[normalizePath(path)] = hashString
	if err := writeIndex(indexPath, entries); err != nil {
		return err
	}

	fmt.Println(hashString)
	// fmt.Println(path)
	// fmt.Println(fileInfo.Name())
	// fmt.Println(fileInfo.ModTime())
	// fmt.Println(fileInfo.Size())
	// fmt.Println(fileInfo.IsDir())
	return nil
}
