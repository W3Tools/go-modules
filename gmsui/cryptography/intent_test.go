package cryptography

import (
	"bytes"
	"testing"
)

func TestIntentWithScope(t *testing.T) {
	tests := []struct {
		name     string
		scope    IntentScope
		expected Intent
	}{
		{
			name:     "TransactionData scope",
			scope:    TransactionData,
			expected: Intent{TransactionData, V0, Sui},
		},
		{
			name:     "TransactionEffects scope",
			scope:    TransactionEffects,
			expected: Intent{TransactionEffects, V0, Sui},
		},
		{
			name:     "CheckpointSummary scope",
			scope:    CheckpointSummary,
			expected: Intent{CheckpointSummary, V0, Sui},
		},
		{
			name:     "PersonalMessage scope",
			scope:    PersonalMessage,
			expected: Intent{PersonalMessage, V0, Sui},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntentWithScope(tt.scope)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestMessageWithIntent(t *testing.T) {
	tests := []struct {
		name     string
		scope    IntentScope
		message  []byte
		expected []byte
	}{
		{
			name:     "TransactionData with message",
			scope:    TransactionData,
			message:  []byte("test message"),
			expected: append(IntentWithScope(TransactionData), []byte("test message")...),
		},
		{
			name:     "TransactionEffects with message",
			scope:    TransactionEffects,
			message:  []byte("another message"),
			expected: append(IntentWithScope(TransactionEffects), []byte("another message")...),
		},
		{
			name:     "Empty message",
			scope:    PersonalMessage,
			message:  []byte(""),
			expected: IntentWithScope(PersonalMessage),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MessageWithIntent(tt.scope, tt.message)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
