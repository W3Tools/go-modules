package gmar

import (
	"fmt"
	"net/url"
	"path"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

func (a *ArweaveClient) UploadTo(data []byte, manifest bool) (*types.Transaction, error) {
	tx, err := a.GetTransaction(data, manifest)
	if err != nil {
		return nil, err
	}

	uploder, err := goar.CreateUploader(a.Client, tx, nil)
	if err != nil {
		return nil, fmt.Errorf("goar.CreateUploader err %v", err)
	}

	err = uploder.Once()
	if err != nil {
		return nil, fmt.Errorf("uploder.Once err %v", err)
	}

	return tx, nil
}

func (a *ArweaveClient) UploadToTxHash(data []byte, manifest bool) (string, error) {
	tx, err := a.UploadTo(data, manifest)
	if err != nil {
		return "", err
	}

	return tx.ID, nil
}

func (a *ArweaveClient) UploadToUrl(data []byte, manifest bool) (string, error) {
	tx, err := a.UploadTo(data, manifest)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(a.Node)
	if err != nil {
		return "", fmt.Errorf("url.Parse err %v", err)
	}
	u.Path = path.Join(u.Path, tx.ID)

	return u.String(), nil
}
