package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func Sha1(path string) (string, error) {
	var result string
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}

	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	hashInBytes := hash.Sum(nil)[:20]
	result = hex.EncodeToString(hashInBytes)
	return result, nil
}
