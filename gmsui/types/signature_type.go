package types

import (
	"encoding/json"
	"errors"
)

type Signature interface {
	isSignature()
}

type SignatureEd25519 struct {
	Ed25519SuiSignature string `json:"Ed25519SuiSignature"`
}

type SignatureSecp256k1 struct {
	Secp256k1SuiSignature string `json:"Secp256k1SuiSignature"`
}

type SignatureSecp256r1 struct {
	Secp256r1SuiSignature string `json:"Secp256r1SuiSignature"`
}

func (SignatureEd25519) isSignature()   {}
func (SignatureSecp256k1) isSignature() {}
func (SignatureSecp256r1) isSignature() {}

type SignatureWrapper struct {
	Signature
}

func (w *SignatureWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["Ed25519SuiSignature"] != nil:
		var s SignatureEd25519
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.Signature = s
	case obj["Secp256k1SuiSignature"] != nil:
		var s SignatureSecp256k1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.Signature = s
	case obj["Secp256r1SuiSignature"] != nil:
		var s SignatureSecp256r1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.Signature = s
	default:
		return errors.New("unknown Signature type")
	}

	return nil
}

func (w *SignatureWrapper) MarshalJSON() ([]byte, error) {
	switch s := w.Signature.(type) {
	case SignatureEd25519:
		return json.Marshal(SignatureEd25519{
			Ed25519SuiSignature: s.Ed25519SuiSignature,
		})
	case SignatureSecp256k1:
		return json.Marshal(SignatureSecp256k1{
			Secp256k1SuiSignature: s.Secp256k1SuiSignature,
		})
	case SignatureSecp256r1:
		return json.Marshal(SignatureSecp256r1{
			Secp256r1SuiSignature: s.Secp256r1SuiSignature,
		})
	default:
		return nil, errors.New("unknown Signature type")
	}
}

// --------------CompressedSignature--------------
type CompressedSignature interface {
	isCompressedSignature()
}

type CompressedSignatureEd25519 struct {
	Ed25519 string `json:"Ed25519"`
}

type CompressedSignatureSecp256k1 struct {
	Secp256k1 string `json:"Secp256k1"`
}

type CompressedSignatureSecp256r1 struct {
	Secp256r1 string `json:"Secp256r1"`
}

func (CompressedSignatureEd25519) isCompressedSignature()   {}
func (CompressedSignatureSecp256k1) isCompressedSignature() {}
func (CompressedSignatureSecp256r1) isCompressedSignature() {}

type CompressedSignatureWrapper struct {
	CompressedSignature
}

func (w *CompressedSignatureWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["Ed25519"] != nil:
		var s CompressedSignatureEd25519
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.CompressedSignature = s
	case obj["Secp256k1"] != nil:
		var s CompressedSignatureSecp256k1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.CompressedSignature = s
	case obj["Secp256r1"] != nil:
		var s CompressedSignatureSecp256r1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.CompressedSignature = s
	default:
		return errors.New("unknown CompressedSignature type")
	}

	return nil
}

func (w *CompressedSignatureWrapper) MarshalJSON() ([]byte, error) {
	switch s := w.CompressedSignature.(type) {
	case CompressedSignatureEd25519:
		return json.Marshal(CompressedSignatureEd25519{
			Ed25519: s.Ed25519,
		})
	case CompressedSignatureSecp256k1:
		return json.Marshal(CompressedSignatureSecp256k1{
			Secp256k1: s.Secp256k1,
		})
	case CompressedSignatureSecp256r1:
		return json.Marshal(CompressedSignatureSecp256r1{
			Secp256r1: s.Secp256r1,
		})
	default:
		return nil, errors.New("unknown CompressedSignature type")
	}
}
