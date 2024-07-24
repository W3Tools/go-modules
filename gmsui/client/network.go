package client

func GetFullNodeURL(network string) string {
	switch network {
	case "mainnet":
		return "https://fullnode.mainnet.sui.io:443"
	case "testnet":
		return "https://fullnode.testnet.sui.io:443"
	case "devnet":
		return "https://fullnode.devnet.sui.io:443"
	case "localnet":
		return "http://127.0.0.1:9000"
	default:
		return "https://fullnode.devnet.sui.io:443"
	}
}
