package gm

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"slices"
)

func ReadFileBytes(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, stat.Size())

	_, err = f.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

/*
Get the hash value of the file
The effect is equivalent to the command line sha256sum

	fileBytes, err := gm.ReadFileBytes("example.go")
	if err != nil {
		fmt.Printf("read file bytes err, msg: %v\n", err)
		return
	}

	hash, err := gm.ReadFileHash(fileBytes)
	if err != nil {
		fmt.Printf("read file hash err, msg: %v\n", err)
		return
	}
	fmt.Printf("hash: %v\n", hash)

	-----------------------------------------------------
	FILE_HASH=$(sha256sum example.go)

	hash == FILE_HASH
*/
func ReadFileHash(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

/*
Hash an array into a nested array

	numArray := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	splitArray := gm.HashArrayToArrays(numArray, 3)
	fmt.Printf("new array: %v\n", splitArray)

	Output: [[1 4 7 9] [2 5 10] [3 6 8]]
*/
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
Truncate the string, keep the beginning and ends, and add an ellipsis in the middle to truncate the string
*/
func TruncateString(v string, start, end int) string {
	if len(v) <= start+end {
		return v
	}

	return fmt.Sprintf("%s...%s", v[:start], v[len(v)-end:])
}

/*
Only add the element if it is not already in the slice
*/
func UniqueAppend[S ~[]E, E comparable](s S, v E) S {
	if !slices.Contains(s, v) {
		s = append(s, v)
	}
	return s
}

func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func NewNumber[T number](v T) T {
	return v
}

func NewUint64[T number](v T) uint64 {
	return uint64(v)
}

func NewStringPtr(v string) *string {
	return &v
}
