package b64_test

import (
	"testing"

	"github.com/W3Tools/go-modules/gmsui/b64"
)

func TestFromBase64(t *testing.T) {
	tests := []struct {
		name          string
		base64String  string
		expectedBytes []byte
		expectError   bool
	}{
		{
			name:          "Valid base64 string",
			base64String:  "SGVsbG8sIFdvcmxkIQ==", // "Hello, World!"
			expectedBytes: []byte("Hello, World!"),
			expectError:   false,
		},
		{
			name:          "Invalid base64 string",
			base64String:  "InvalidBase64",
			expectedBytes: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := b64.FromBase64(tt.base64String)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && string(result) != string(tt.expectedBytes) {
				t.Errorf("expected %v, but got %v", tt.expectedBytes, result)
			}
		})
	}
}

func TestToBase64(t *testing.T) {
	tests := []struct {
		name           string
		inputBytes     []byte
		expectedBase64 string
	}{
		{
			name:           "Encode string to base64",
			inputBytes:     []byte("Hello, World!"),
			expectedBase64: "SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:           "Encode empty string to base64",
			inputBytes:     []byte(""),
			expectedBase64: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b64.ToBase64(tt.inputBytes)
			if result != tt.expectedBase64 {
				t.Errorf("expected %v, but got %v", tt.expectedBase64, result)
			}
		})
	}
}
