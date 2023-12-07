package gmsui

import "github.com/block-vision/sui-go-sdk/sui"

type SuiNetwork = string

var (
	MainnetFullNode SuiNetwork = "https://fullnode.mainnet.sui.io:443/"
	TestnetFullNode SuiNetwork = "https://fullnode.testnet.sui.io:443/"
	DevnetFullNode  SuiNetwork = "https://fullnode.devnet.sui.io:443/"
)

// Create New Provider
func NewSuiProviderFromNetwork(network SuiNetwork) sui.ISuiAPI {
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

func NewSuiProvider(rpc string) sui.ISuiAPI {
	return sui.NewSuiClient(rpc)
}
