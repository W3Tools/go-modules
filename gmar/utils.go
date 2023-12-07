package gmar

import (
	"fmt"
	"math/big"
	"net/http"

	gm "github.com/W3Tools/go-modules"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

func (a *ArweaveClient) GetBalance() (arAmount *big.Float, err error) {
	return a.Client.GetWalletBalance(a.Wallet.Signer.Address)
}

func (a *ArweaveClient) GetTxPrice(data []byte) (int64, error) {
	return a.Client.GetTransactionPrice(data, nil)
}

func (a *ArweaveClient) GetTransaction(data []byte, manifest bool) (*types.Transaction, error) {
	anchor, err := a.Client.GetTransactionAnchor()
	if err != nil {
		return nil, fmt.Errorf("client.GetTransactionAnchor err %v", err)
	}

	reward, err := a.GetTxPrice(data)
	if err != nil {
		return nil, fmt.Errorf("a.GetTxPrice err %v", err)
	}

	var tags []types.Tag

	fileHash, _ := gm.ReadFileHash(data)

	if manifest {
		tags = append(tags, types.Tag{
			Name:  "Content-Type",
			Value: "application/x.arweave-manifest+json",
		})
	} else {
		tags = append(tags, types.Tag{
			Name:  "Content-Type",
			Value: http.DetectContentType(data),
		})
		tags = append(tags, types.Tag{
			Name:  "User-Agent",
			Value: "W3Tools",
		})
		tags = append(tags, types.Tag{
			Name:  "FileHash",
			Value: fileHash,
		})
	}

	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     utils.TagsEncode(tags),
		Data:     utils.Base64Encode(data),
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100)/100),
		LastTx:   anchor,
		Owner:    utils.Base64Encode(a.Wallet.Signer.PubKey.N.Bytes()),
	}

	err = utils.SignTransaction(tx, a.Wallet.Signer.PrvKey)
	if err != nil {
		return nil, fmt.Errorf("utils.SignTransaction err %v", err)
	}

	return tx, nil
}
