package transactions

import (
	"fmt"
	"strconv"

	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
)

type ObjectStringRef struct {
	ObjectId string `json:"objectId"`
	Version  string `json:"version"`
	Digest   string `json:"digest"`
}

func (ref ObjectStringRef) toObjectRef() (*sui_types.ObjectRef, error) {
	objectId, err := sui_types.NewObjectIdFromHex(ref.ObjectId)
	if err != nil {
		return nil, fmt.Errorf("can not create object id from hex [%s]: %v", ref.ObjectId, err)
	}
	digest, err := sui_types.NewDigest(ref.Digest)
	if err != nil {
		return nil, fmt.Errorf("can not create digest [%s]: %v", ref.Digest, err)
	}
	version, err := strconv.ParseUint(ref.Version, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can not parse version [%s] to uint64: %v", ref.Version, err)
	}

	return &sui_types.ObjectRef{ObjectId: *objectId, Version: version, Digest: *digest}, nil
}

func coinStructToObjectRef(coin types.CoinStruct) (*sui_types.ObjectRef, error) {
	return ObjectStringRef{ObjectId: coin.CoinObjectId, Version: coin.Version, Digest: coin.Digest}.toObjectRef()
}
