package cryptography_test

import (
	"bytes"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
)

func TestIntentWithScope(t *testing.T) {
	tests := []struct {
		name     string
		scope    cryptography.IntentScope
		expected cryptography.Intent
	}{
		{
			name:     "TransactionData scope",
			scope:    cryptography.TransactionData,
			expected: cryptography.Intent{cryptography.TransactionData, cryptography.V0, cryptography.Sui},
		},
		{
			name:     "TransactionEffects scope",
			scope:    cryptography.TransactionEffects,
			expected: cryptography.Intent{cryptography.TransactionEffects, cryptography.V0, cryptography.Sui},
		},
		{
			name:     "CheckpointSummary scope",
			scope:    cryptography.CheckpointSummary,
			expected: cryptography.Intent{cryptography.CheckpointSummary, cryptography.V0, cryptography.Sui},
		},
		{
			name:     "PersonalMessage scope",
			scope:    cryptography.PersonalMessage,
			expected: cryptography.Intent{cryptography.PersonalMessage, cryptography.V0, cryptography.Sui},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cryptography.IntentWithScope(tt.scope)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestMessageWithIntent(t *testing.T) {
	tests := []struct {
		name     string
		scope    cryptography.IntentScope
		message  []byte
		expected []byte
	}{
		{
			name:     "TransactionData with message",
			scope:    cryptography.TransactionData,
			message:  []byte("test message"),
			expected: append(cryptography.IntentWithScope(cryptography.TransactionData), []byte("test message")...),
		},
		{
			name:     "TransactionEffects with message",
			scope:    cryptography.TransactionEffects,
			message:  []byte("another message"),
			expected: append(cryptography.IntentWithScope(cryptography.TransactionEffects), []byte("another message")...),
		},
		{
			name:     "Empty message",
			scope:    cryptography.PersonalMessage,
			message:  []byte(""),
			expected: cryptography.IntentWithScope(cryptography.PersonalMessage),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cryptography.MessageWithIntent(tt.scope, tt.message)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
