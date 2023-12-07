package gm

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func ReadFileBytes(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open err, msg: %v", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("f.Stat err, msg: %v", err)
	}

	buffer := make([]byte, stat.Size())

	_, err = f.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("f.Read err, msg: %v", err)
	}
	return buffer, nil
}

func ReadFileHash(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", fmt.Errorf("io.Copy err %v", err)
	}

	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}
