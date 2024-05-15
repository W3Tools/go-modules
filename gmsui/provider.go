package gmsui

import (
	sdk_client "github.com/W3Tools/go-sui-sdk/v2/client"
)

type SuiNetwork = string

var (
	MainnetFullNode SuiNetwork = "https://fullnode.mainnet.sui.io:443/"
	TestnetFullNode SuiNetwork = "https://fullnode.testnet.sui.io:443/"
	DevnetFullNode  SuiNetwork = "https://fullnode.devnet.sui.io:443/"
)

// Create New Provider
func NewSuiProviderFromNetwork(network SuiNetwork) (*sdk_client.Client, error) {
	switch network {
	case "mainnet":
		return NewSuiProvider(MainnetFullNode)
	case "testnet":
		return NewSuiProvider(TestnetFullNode)
	case "devnet":
		return NewSuiProvider(DevnetFullNode)
	default:
		return NewSuiProvider(DevnetFullNode)
	}
}

func NewSuiProvider(rpc string) (*sdk_client.Client, error) {
	return sdk_client.Dial(rpc)
}
