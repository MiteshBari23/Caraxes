package repo

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func hashBlob(content []byte) (string,[]byte, error) {
	header := fmt.Sprintf("blob %d\x00", len(content))

	hasher := sha1.New()
	data := append([]byte(header), content...)
	hasher.Write(data)

	sum := hasher.Sum(nil)
	hashString := hex.EncodeToString(sum)

	return hashString, data, nil
}
