package repo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readCommitBody(objectDir string, hash string) (string, error) {
	data, err := readObject(objectDir, hash)
	if err != nil {
		return "", err
	}

	nullIndex := bytes.IndexByte(data, 0)
	body := string(data[nullIndex+1:])
	return body, nil
}

type CommitInfo struct {
	Tree    string
	Parent  string
	Author  string
	Message string
}

func parseCommit(body string) CommitInfo {
	info := CommitInfo{}
	lines := strings.Split(body, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "tree ") {
			info.Tree = strings.TrimPrefix(line, "tree ")
		} else if strings.HasPrefix(line, "parent ") {
			info.Parent = strings.TrimPrefix(line, "parent ")
		} else if strings.HasPrefix(line, "author ") {
			info.Author = strings.TrimPrefix(line, "author ")
		}
	}

	parts := strings.SplitN(body, "\n\n", 2)
	if len(parts) == 2 {
		info.Message = parts[1]
	}

	return info
}

func LogRepo() error {
	objectDir := filepath.Join(".caraxes", "objects")
	refPath := filepath.Join(".caraxes", "refs", "heads", "master")

	data, err := os.ReadFile(refPath)
	if err != nil {
		return fmt.Errorf("no commits yet")
	}
	currentHash := strings.TrimSpace(string(data))

	for currentHash != "" {
		body, err := readCommitBody(objectDir, currentHash)
		if err != nil {
			return err
		}
		info := parseCommit(body)

		fmt.Printf("commit %s\n", currentHash)
		fmt.Printf("Author: %s\n", info.Author)
		fmt.Printf("\n    %s\n\n", strings.TrimSpace(info.Message))

		currentHash = info.Parent
	}

	return nil
}
