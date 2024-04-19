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

func HashArrayToArrays[T any](input []T, count int) [][]T {
	result := make([][]T, count)
	for i := range result {
		result[i] = make([]T, 0)
	}

	avg := len(input) / count
	extra := len(input) % count

	index := 0
	for _, val := range input {
		result[index] = append(result[index], val)
		if len(result[index]) == avg && extra > 0 {
			extra--
			index++
		}
		index++
		index %= count
	}

	return result
}

func Map[T any, T2 any](s []T, fn func(t T) (T2, error)) (t2 []T2, err error) {
	for i := range s {
		v, err := fn(s[i])
		if err != nil {
			return t2, err
		}
		t2 = append(t2, v)
	}
	return
}
