package gmar

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/everFinance/goar"
)

type ArweaveClient struct {
	Node   string
	Wallet *goar.Wallet
	Client *goar.Client
}

const (
	ManifestFile        = "manifest.w3tools"
	ManifestContentType = "application/x.arweave-manifest+json"
	IndexFile           = "index.html"
)

func InitArweaveClient(keyFile, node string) (*ArweaveClient, error) {
	if strings.EqualFold(node, "") {
		node = "https://arweave.net/"
	}

	wallet, err := goar.NewWalletFromPath(keyFile, node)
	if err != nil {
		return nil, fmt.Errorf("goar.NewWalletFromPath err %v", err)
	}

	_arweave := &ArweaveClient{
		Node:   node,
		Wallet: wallet,
		Client: wallet.Client,
	}
	return _arweave, nil
}

func (a *ArweaveClient) NewManifest(path string) (*ArweaveManifest, error) {
	var manifest = &ArweaveManifest{}
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		manifest = NewManifest()

		data, err := json.Marshal(manifest)
		if err != nil {
			return nil, err
		}

		err = WriteManifest(path, data)
		if err != nil {
			return nil, err
		}
		return manifest, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if manifest.Paths == nil {
		manifest.Paths = make(map[string]ArweaveManifestPath)
	}

	err = json.Unmarshal(data, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
