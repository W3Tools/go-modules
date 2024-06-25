package gm

import (
	"crypto/sha256"
	"fmt"
	"os"
	"reflect"
	"slices"
	"testing"
)

func TestReadFileBytes(t *testing.T) {
	// Creating a temporary file
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Unable to create temporary file, msg: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write data to a temporary file
	testData := []byte("Hello, GO Modules!")
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatalf("Unable to write temporary file, msg: %v", err)
	}

	// Close temporary files
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Unable to close temporary file: %v", err)
	}

	// Calling the ReadFileBytes function
	result, err := ReadFileBytes(tempFile.Name())
	if err != nil {
		t.Fatalf("ReadFileBytes function err, msg: %v", err)
	}

	if !reflect.DeepEqual(result, testData) {
		t.Errorf("Expected to read data %v, but got %v", testData, result)
	}
}

func TestReadFileHash(t *testing.T) {
	testData := []byte("Hello, GO Modules!")

	hash := sha256.New()
	if _, err := hash.Write(testData); err != nil {
		t.Fatalf("Unable to compute expected hash value, msg: %v", err)
	}
	expectedHash := fmt.Sprintf("%x", hash.Sum(nil))

	result, err := ReadFileHash(testData)
	if err != nil {
		t.Fatalf("ReadFileHash function err, msg: %v", err)
	}

	if result != expectedHash {
		t.Errorf("Expected to read data %v, but got %v", expectedHash, result)
	}
}

func TestHashArrayToArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		count    int
		expected [][]int
	}{
		{
			name:     "Equal distribution",
			input:    []int{1, 2, 3, 4, 5, 6},
			count:    3,
			expected: [][]int{{1, 4}, {2, 5}, {3, 6}},
		},
		{
			name:     "One bucket",
			input:    []int{1, 2, 3},
			count:    1,
			expected: [][]int{{1, 2, 3}},
		},
		{
			name:     "More buckets than elements",
			input:    []int{1, 2, 3},
			count:    5,
			expected: [][]int{{1}, {2}, {3}, {}, {}},
		},
		{
			name:     "Empty input",
			input:    []int{},
			count:    3,
			expected: [][]int{{}, {}, {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashArrayToArrays(tt.input, tt.count)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{2, 4, 6}
	result, err := Map(input, func(t int) (int, error) {
		return t * 2, nil
	})
	if err != nil {
		t.Fatalf("Map returned an error: %v", err)
	}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, but got %v", expected, result)
	}

	// Test error case
	_, err = Map(input, func(t int) (int, error) {
		if t == 2 {
			return 0, fmt.Errorf("error on 2")
		}
		return t * 2, nil
	})
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func TestFilterOne(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := 3
	result := FilterOne(input, func(t int) bool {
		return t == 3
	})
	if result != expected {
		t.Errorf("expected %v, but got %v", expected, result)
	}

	// Test not found case
	result = FilterOne(input, func(t int) bool {
		return t == 6
	})
	if result != 0 {
		t.Errorf("expected 0, but got %v", result)
	}
}

func TestTruncateString(t *testing.T) {
	input := "Hello, Go Modules!"
	expected := "He...s!"
	result := TruncateString(input, 2, 2)
	if result != expected {
		t.Errorf("expected %v, but got %v", expected, result)
	}

	// Test case where string is not long enough
	expected = "Hello"
	result = TruncateString("Hello", 2, 3)
	if result != expected {
		t.Errorf("expected %v, but got %v", expected, result)
	}
}

func TestUniqueAppend(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{1, 2, 3, 4}
	result := UniqueAppend(input, 4)
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, but got %v", expected, result)
	}

	// Test case where element is already present
	expected = []int{1, 2, 3}
	result = UniqueAppend(input, 3)
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, but got %v", expected, result)
	}
}

func TestBytesEqual(t *testing.T) {
	a := []byte{1, 2, 3}
	b := []byte{1, 2, 3}
	if !BytesEqual(a, b) {
		t.Error("expected true, but got false")
	}

	// Test case where slices have different lengths
	b = []byte{1, 2}
	if BytesEqual(a, b) {
		t.Error("expected false, but got true")
	}

	// Test case where slices have different contents
	b = []byte{1, 2, 4}
	if BytesEqual(a, b) {
		t.Error("expected false, but got true")
	}
}
