package gmsui

import (
	"fmt"
	"math/big"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
)

type VectorU8 struct {
	Data []uint8
}

type VectorU16 struct {
	Data []uint16
}

type VectorU32 struct {
	Data []uint32
}

type VectorU64 struct {
	Data []uint64
}

type VectorBigInt struct {
	Data []big.Int
}

type SuiAddress = string

type VectorAddress struct {
	Data []SuiAddress
}

func (v *VectorU8) Marshal() (buf []byte, err error) {
	buf = append(buf, byte(len(v.Data)))

	for _, v := range v.Data {
		bcsData, err := bcs.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("bcs.Marshal %v", err)
		}
		buf = append(buf, bcsData...)
	}
	return
}

func (v *VectorU16) Marshal() (buf []byte, err error) {
	buf = append(buf, byte(len(v.Data)))

	for _, v := range v.Data {
		bcsData, err := bcs.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("bcs.Marshal %v", err)
		}
		buf = append(buf, bcsData...)
	}
	return
}

func (v *VectorU32) Marshal() (buf []byte, err error) {
	buf = append(buf, byte(len(v.Data)))

	for _, v := range v.Data {
		bcsData, err := bcs.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("bcs.Marshal %v", err)
		}
		buf = append(buf, bcsData...)
	}
	return
}

func (v *VectorU64) Marshal() (buf []byte, err error) {
	buf = append(buf, byte(len(v.Data)))

	for _, v := range v.Data {
		bcsData, err := bcs.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("bcs.Marshal %v", err)
		}
		buf = append(buf, bcsData...)
	}
	return
}

func (v *VectorBigInt) Marshal() (buf []byte, err error) {
	buf = append(buf, byte(len(v.Data)))

	for _, v := range v.Data {
		bcsData, err := bcs.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("bcs.Marshal %v", err)
		}
		buf = append(buf, bcsData...)
	}
	return
}

func (v *VectorAddress) Marshal() (buf []byte, err error) {
	buf = append(buf, byte(len(v.Data)))

	for _, v := range v.Data {
		addressHex, err := sui_types.NewAddressFromHex(v)
		if err != nil {
			return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
		}
		bcsData, err := bcs.Marshal(addressHex)
		if err != nil {
			return nil, fmt.Errorf("bcs.Marshal %v", err)
		}
		buf = append(buf, bcsData...)
	}
	return
}
