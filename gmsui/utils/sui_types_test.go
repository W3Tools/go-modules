package utils_test

import (
	"reflect"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/utils"
)

func TestNormalizeSuiAddress(t *testing.T) {
	tests := []struct {
		str      string
		expected string
	}{
		{
			str:      "",
			expected: "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			str:      "a",
			expected: "0x000000000000000000000000000000000000000000000000000000000000000a",
		},
		{
			str:      "123",
			expected: "0x0000000000000000000000000000000000000000000000000000000000000123",
		},
		{
			str:      "0xa123",
			expected: "0x000000000000000000000000000000000000000000000000000000000000a123",
		},
		{
			str:      "0x6",
			expected: "0x0000000000000000000000000000000000000000000000000000000000000006",
		},
		{
			str:      "0x06",
			expected: "0x0000000000000000000000000000000000000000000000000000000000000006",
		},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			normalizedAddress := utils.NormalizeSuiAddress(tt.str)
			if !reflect.DeepEqual(normalizedAddress, tt.expected) {
				t.Errorf("normalize sui address expected %s, but got %v", tt.expected, normalizedAddress)
			}

			normalizedObject := utils.NormalizeSuiObjectId(tt.str)
			if !reflect.DeepEqual(normalizedObject, tt.expected) {
				t.Errorf("normalize sui object expected %s, but got %v", tt.expected, normalizedObject)
			}
		})
	}
}
func TestNormalizeShortSuiAddress(t *testing.T) {
	tests := []struct {
		str      string
		expected string
	}{
		{
			str:      "0x0000000000000000000000000000000000000000000000000000000000000000",
			expected: "0x0",
		},
		{
			str:      "0x000000000000000000000000000000000000000000000000000000000000000a",
			expected: "0xa",
		},
		{
			str:      "0x0000000000000000000000000000000000000000000000000000000000000123",
			expected: "0x123",
		},
		{
			str:      "0x000000000000000000000000000000000000000000000000000000000000a123",
			expected: "0xa123",
		},
		{
			str:      "0x0000000000000000000000000000000000000000000000000000000000000006",
			expected: "0x6",
		},
		{
			str:      "0x0000000000000000000000000000000000000000000000000000000000000006",
			expected: "0x6",
		},
		{
			str:      "0x06",
			expected: "0x6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			normalizedShortSuiAddress := utils.NormalizeShortSuiAddress(tt.str)
			if !reflect.DeepEqual(normalizedShortSuiAddress, tt.expected) {
				t.Errorf("normalize short sui address expected %s, but got %s", tt.expected, normalizedShortSuiAddress)
			}

			normalizedShortSuiObjectId := utils.NormalizeShortSuiObjectId(tt.str)
			if !reflect.DeepEqual(normalizedShortSuiObjectId, tt.expected) {
				t.Errorf("normalize short sui object id expected %s, but got %s", tt.expected, normalizedShortSuiObjectId)
			}
		})
	}
}

func TestNormalizeSuiCoinType(t *testing.T) {
	tests := []struct {
		str      string
		expected string
	}{
		{
			str:      "0x2::sui::SUI",
			expected: "0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI",
		},
		{
			str:      "0x02::sui::SUI",
			expected: "0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI",
		},
		{
			str:      "0x01",
			expected: "0x01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			normalizedShortSuiAddress := utils.NormalizeSuiCoinType(tt.str)
			if !reflect.DeepEqual(normalizedShortSuiAddress, tt.expected) {
				t.Errorf("normalize short sui address expected %s, but got %s", tt.expected, normalizedShortSuiAddress)
			}
		})
	}
}

func TestNormalizeShortSuiCoinType(t *testing.T) {
	tests := []struct {
		str      string
		expected string
	}{
		{
			str:      "0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI",
			expected: "0x2::sui::SUI",
		},
		{
			str:      "0x02::sui::SUI",
			expected: "0x2::sui::SUI",
		},
		{
			str:      "0x01",
			expected: "0x01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			normalizedShortSuiAddress := utils.NormalizeShortSuiCoinType(tt.str)
			if !reflect.DeepEqual(normalizedShortSuiAddress, tt.expected) {
				t.Errorf("normalize short sui address expected %s, but got %s", tt.expected, normalizedShortSuiAddress)
			}
		})
	}
}
