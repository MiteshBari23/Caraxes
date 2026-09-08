package repo

import (
	"bytes"
	"compress/zlib"
	"os"
	"path/filepath"
)

func writeObject(objectDir string, hash string, data []byte) error {
	dir := filepath.Join(objectDir, hash[:2])  // :2 means the folder name will be first to char of hash
	objectPath := filepath.Join(dir, hash[2:]) // 2: 2nd hash char ke baad ke sare char milake filename banaynege

	if _, err := os.Stat(objectPath); err == nil {
		return nil
	}

	//0755: Owner can read/write/execute; group and others can read/execute (standard for directories and executable programs).
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var buf bytes.Buffer

	writer := zlib.NewWriter(&buf)
	writer.Write(data)
	writer.Close() // it is necessary or else all compressed data will not be flushed in buf

	//0644: Owner can read/write; group and others can read (standard for regular files).
	return os.WriteFile(objectPath, buf.Bytes(), 0644)

}
