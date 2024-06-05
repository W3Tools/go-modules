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

/*
Calls a defined callback function on each element of an array, and returns an array that contains the results.

	func StringArrayToNumberArray() {
		var stringArray []string = []string{"1", "2", "3"}

		numberArray, err := Map(stringArray, func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		})
		if err != nil {
			fmt.Println("string array to number array error, msg: ", err)
			return
		}
		fmt.Println(numberArray)
	}

	func CalculateSquareValue() {
		var numArray []int64 = []int64{2, 3, 4}
		squareArray, _ := Map(numArray, func(n int64) (int64, error) {
			return (n * n), nil
		})
		fmt.Println(squareArray)
	}
*/
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

/*
Returns the first element that meet the condition specified in a callback function.
*/
func FilterOne[T any](s []T, fn func(T) bool) (t T) {
	for i := range s {
		if fn(s[i]) {
			return s[i]
		}
	}
	return
}

/*
Cut off the beginning and end of the string and add ellipsis at both ends to truncate the string
*/
func TruncateString(v string, start, end int) string {
	if len(v) < start+end {
		return v
	}

	return fmt.Sprintf("%s...%s", v[:start], v[len(v)-end:])
}
