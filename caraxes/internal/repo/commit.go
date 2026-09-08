package repo

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func writeCommitTree(objectDir string, treeHash string, parentHash string, message string) (string, error) {
	authorEmail := os.Getenv("XGIT_AUTHOR_EMAIL")
	authorName := os.Getenv("XGIT_AUTHOR_NAME")

	if authorName == "" || authorEmail == "" {
		return "", fmt.Errorf("author identity missing: set XGIT_AUTHOR_NAME and XGIT_AUTHOR_EMAIL")
	}
	now := time.Now()
	timeStamp := fmt.Sprintf("%d %s", now.Unix(), now.Format("-0700"))

	var content strings.Builder

	content.WriteString(fmt.Sprintf("tree %s\n", treeHash))

	if parentHash != "" {
		content.WriteString(fmt.Sprintf("parent %s\n", parentHash))
	}

	content.WriteString(fmt.Sprintf("author %s <%s> %s\n", authorName, authorEmail, timeStamp))
	content.WriteString(fmt.Sprintf("committer %s <%s> %s\n", authorName, authorEmail, timeStamp))

	content.WriteString("\n")
	content.WriteString(message)
	content.WriteString("\n")

	contentBytes := []byte(content.String())
	header := fmt.Sprintf("commit %d\x00", len(contentBytes))
	data := append([]byte(header), contentBytes...)

	sum := sha1.Sum(data)
	hashString := hex.EncodeToString(sum[:])

	if err := writeObject(objectDir, hashString, data); err != nil {
		return "", err
	}

	return hashString, nil
}

func CommitRepo(message string) error {
	objectDir := filepath.Join(".caraxes", "objects")
	indexPath := filepath.Join(".caraxes", "index")
	refPath := filepath.Join(".caraxes", "refs", "heads", "master")

	entries, err := readIndex(indexPath)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("nothing to commit, working tree clean")
	}

	root := buildTree(entries)
	treeHash, err := writeTree(objectDir, root)
	if err != nil {
		return err
	}

	var parentHash string
	if data, err := os.ReadFile(refPath); err == nil {
		parentHash = strings.TrimSpace(string(data))
	}

	commitHash, err := writeCommitTree(objectDir, treeHash, parentHash, message)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(refPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(refPath, []byte(commitHash+"\n"), 0644); err != nil {
		return err
	}

	fmt.Printf("committed as %s\n", commitHash)

	return nil
}
