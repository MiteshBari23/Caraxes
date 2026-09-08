package repo

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readIndex(indexPath string) (map[string]string, error) {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return map[string]string{}, nil
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	entries := map[string]string{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
func writeIndex(indexPath string, entries map[string]string) error {
	data, err := json.MarshalIndent(entries, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(indexPath, data, 0644)
}

func normalizePath(path string) string {
	return filepath.ToSlash(path)
}


