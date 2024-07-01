package cryptography

import (
	"fmt"
	"io"

	"github.com/fardream/go-bcs/bcs"
)

type CompressedSignature struct {
	Signature [65]byte `json:"signature"`
}

type PubKeyEnumWeightPair struct {
	PubKey []byte `json:"pubKey"`
	Weight uint8  `json:"weight"`
}

type MultiSigPublicKeyStruct struct {
	PubKeyMap []*PubKeyEnumWeightPair `json:"pubKeymap"`
	Threshold uint16                  `json:"threshold"`
}

type MultiSigStruct struct {
	Sigs           []CompressedSignature   `json:"sigs"`
	Bitmap         uint16                  `json:"bitmap"`
	MultisigPubKey MultiSigPublicKeyStruct `json:"multisigPubKey"`
}

type MultiSigPublicKeyPair struct {
	Weight    uint8     `json:"weight"`
	PublicKey PublicKey `json:"publicKey"`
}

var (
	_ bcs.Marshaler   = (*PubKeyEnumWeightPair)(nil)
	_ bcs.Unmarshaler = (*PubKeyEnumWeightPair)(nil)
)

func (p PubKeyEnumWeightPair) MarshalBCS() ([]byte, error) {
	switch len(p.PubKey) {
	case 33:
		data := struct {
			PubKey [33]byte `json:"pubKey"`
			Weight uint8    `json:"weight"`
		}{PubKey: [33]byte(p.PubKey[:]), Weight: p.Weight}
		return bcs.Marshal(data)
	case 34:
		data := struct {
			PubKey [34]byte `json:"pubKey"`
			Weight uint8    `json:"weight"`
		}{PubKey: [34]byte(p.PubKey[:]), Weight: p.Weight}
		return bcs.Marshal(data)
	default:
		return nil, fmt.Errorf("pubKey length must be either 33 or 34 bytes")

	}
}

func (p *PubKeyEnumWeightPair) UnmarshalBCS(r io.Reader) (int, error) {
	pubkeyType, n, err := ReadByte(r)
	if err != nil {
		return n, err
	}

	pubkeySize := 0
	switch pubkeyType {
	case SignatureSchemeToFlag[Ed25519Scheme]:
		pubkeySize = Ed25519PublicKeySize
	case SignatureSchemeToFlag[Secp256k1Scheme]:
		pubkeySize = Secp256k1PublicKeySize
	case SignatureSchemeToFlag[Secp256r1Scheme]:
		pubkeySize = Secp256r1PublicKeySize
	}

	pubkey := make([]byte, pubkeySize)
	read, err := r.Read(pubkey)
	n += read
	if err != nil {
		return n, err
	}

	weight, read, err := ReadByte(r)
	n += read
	if err != nil {
		return n, err
	}

	p.PubKey = append([]byte{pubkeyType}, pubkey...)
	p.Weight = weight

	return n, nil
}

func ReadByte(r io.Reader) (byte, int, error) {
	b := [1]byte{}
	n, err := r.Read(b[:])
	if err != nil {
		return 0, n, err
	}

	if n == 0 {
		return 0, n, io.ErrUnexpectedEOF
	}

	return b[0], n, nil
}
