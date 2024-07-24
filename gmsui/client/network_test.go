package client_test

import (
	"testing"

	"github.com/W3Tools/go-modules/gmsui/client"
)

func TestGetFullNodeURL(t *testing.T) {
	tests := []struct {
		network string
		want    string
	}{
		{"mainnet", "https://fullnode.mainnet.sui.io:443"},
		{"testnet", "https://fullnode.testnet.sui.io:443"},
		{"devnet", "https://fullnode.devnet.sui.io:443"},
		{"localnet", "http://127.0.0.1:9000"},
		{"unknown", "https://fullnode.devnet.sui.io:443"},
	}

	for _, tt := range tests {
		t.Run(tt.network, func(t *testing.T) {
			got := client.GetFullNodeURL(tt.network)
			if got != tt.want {
				t.Errorf("GetFullNodeURL(%q) = %v; want %v", tt.network, got, tt.want)
			}
		})
	}
}
