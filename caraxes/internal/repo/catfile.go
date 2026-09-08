package repo

import (
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
)

func readObject(objectDir string, hash string) ([]byte, error) {

	dir := filepath.Join(objectDir, hash[:2])
	objectPath := filepath.Join(dir, hash[2:])

	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	content, err := io.ReadAll(zlibReader)
	if err != nil {
		return nil, err
	}

	return content, nil

}


