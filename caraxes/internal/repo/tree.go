package repo

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type TreeNode struct {
	isFile   bool
	hash     string
	children map[string]*TreeNode
}

func buildTree(entries map[string]string) *TreeNode {
	root := &TreeNode{isFile: false, children: map[string]*TreeNode{}}

	for path, hash := range entries {
		parts := strings.Split(path, "/")
		current := root

		for _, part := range parts[:len(parts)-1] {
			if current.children[part] == nil {
				current.children[part] = &TreeNode{isFile: false, children: map[string]*TreeNode{}}
			}
			current = current.children[part]
		}
		fileName := parts[len(parts)-1]
		current.children[fileName] = &TreeNode{isFile: true, hash: hash}
	}

	return root
}

func writeTree(objectsDir string, node *TreeNode) (string, error) {
	names := make([]string, 0, len(node.children))
	for name := range node.children {
		names = append(names, name)
	}
	sort.Strings(names)

	var entries bytes.Buffer

	for _, name := range names {
		child := node.children[name]

		var mode string
		var hashHex string

		if child.isFile {
			mode = "100644"
			hashHex = child.hash
		} else {
			mode = "40000"
			subHash, err := writeTree(objectsDir, child)
			if err != nil {
				return "", err
			}
			hashHex = subHash
		}
		hashBytes, err := hex.DecodeString(hashHex)
		if err != nil {
			return "", err
		}

		entries.WriteString(mode + " " + name)
		entries.WriteByte(0)
		entries.Write(hashBytes)

	}
	header := fmt.Sprintf("tree %d\x00", entries.Len())
	data := append([]byte(header), entries.Bytes()...)

	sum := sha1.Sum(data)
	hashString := hex.EncodeToString(sum[:])

	if err := writeObject(objectsDir, hashString, data); err != nil {
		return "", err
	}
	return hashString, nil

}
