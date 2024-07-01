package multisig

import (
	"bytes"
	"encoding/hex"
	"fmt"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/W3Tools/go-modules/gmsui/utils"
	"github.com/W3Tools/go-modules/gmsui/verify"
	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

const (
	MaxSignerInMultisig = 10
	MinSignerInMultisig = 1
)

var (
	_ cryptography.PublicKey = (*MultiSigPublicKey)(nil)
)

type StringPubKeyEnumWeightPair struct {
	PubKey string `json:"pubKey"`
	Weight uint8  `json:"weight"`
}

type ParsedPartialMultiSigSignature struct {
	SignatureScheme cryptography.SignatureScheme
	Signature       []byte
	PublicKey       cryptography.PublicKey
	Weight          uint8
}

type MultiSigPublicKey struct {
	rawBytes          []byte
	multisigPublicKey cryptography.MultiSigPublicKeyStruct
	publicKeys        []cryptography.MultiSigPublicKeyPair
	cryptography.BasePublicKey
}

type PublicKeyWeightPair struct {
	PublicKey cryptography.PublicKey
	Weight    uint8
}

func NewMultiSigPublicKey[T string | []byte | cryptography.MultiSigPublicKeyStruct](value T) (multisig *MultiSigPublicKey, err error) {
	multisig = new(MultiSigPublicKey)
	switch v := any(value).(type) {
	case string:
		rawBytes, err := b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
		multisig.rawBytes = rawBytes

		multisigPublicKeyStruct := new(cryptography.MultiSigPublicKeyStruct)
		_, err = bcs.Unmarshal(multisig.rawBytes, &multisigPublicKeyStruct)
		if err != nil {
			return nil, err
		}
		multisig.multisigPublicKey = *multisigPublicKeyStruct
	case []byte:
		multisig.rawBytes = v

		multisigPublicKeyStruct := new(cryptography.MultiSigPublicKeyStruct)
		_, err = bcs.Unmarshal(multisig.rawBytes, &multisigPublicKeyStruct)
		if err != nil {
			return nil, err
		}
		multisig.multisigPublicKey = *multisigPublicKeyStruct
	case cryptography.MultiSigPublicKeyStruct:
		multisig.multisigPublicKey = v

		rawBytes, err := bcs.Marshal(multisig.multisigPublicKey)
		if err != nil {
			return nil, err
		}
		multisig.rawBytes = rawBytes
	}
	if multisig.multisigPublicKey.Threshold < 1 {
		return nil, fmt.Errorf("invalid threshold")
	}

	seenPublicKeys := make(map[string]bool)
	multisig.publicKeys, err = gm.Map(multisig.multisigPublicKey.PubKeyMap, func(v *cryptography.PubKeyEnumWeightPair) (cryptography.MultiSigPublicKeyPair, error) {
		publicKeyString := string(v.PubKey[:])
		if ok := seenPublicKeys[publicKeyString]; ok {
			return cryptography.MultiSigPublicKeyPair{}, fmt.Errorf("multisig does not support duplicate public keys")
		}
		seenPublicKeys[publicKeyString] = true

		if v.Weight < 1 {
			return cryptography.MultiSigPublicKeyPair{}, fmt.Errorf("invalid weight")
		}

		pubKey, err := verify.PublicKeyFromRawBytes(cryptography.SignatureFlagToScheme[v.PubKey[0]], v.PubKey[1:])
		if err != nil {
			return cryptography.MultiSigPublicKeyPair{}, err
		}
		return cryptography.MultiSigPublicKeyPair{PublicKey: pubKey, Weight: v.Weight}, nil
	})
	if err != nil {
		return nil, err
	}

	var totalWeight uint16 = 0
	gm.Map(multisig.publicKeys, func(v cryptography.MultiSigPublicKeyPair) (any, error) {
		totalWeight = totalWeight + uint16(v.Weight)
		return nil, nil
	})

	if multisig.multisigPublicKey.Threshold > totalWeight {
		return nil, fmt.Errorf("unreachable threshold")
	}

	if len(multisig.publicKeys) > MaxSignerInMultisig {
		return nil, fmt.Errorf("max number of signers in a multisig is %d", MaxSignerInMultisig)
	}

	if len(multisig.publicKeys) < MinSignerInMultisig {
		return nil, fmt.Errorf("min number of signers in a multisig is %d", MinSignerInMultisig)
	}

	multisig.SetSelf(multisig)
	return
}

func (multisig *MultiSigPublicKey) FromPublicKeys(publicKeys []PublicKeyWeightPair, threshold uint16) (*MultiSigPublicKey, error) {
	pubKeyMap, err := gm.Map(publicKeys, func(v PublicKeyWeightPair) (*cryptography.PubKeyEnumWeightPair, error) {
		return &cryptography.PubKeyEnumWeightPair{PubKey: v.PublicKey.ToSuiBytes(), Weight: v.Weight}, nil
	})
	if err != nil {
		return nil, err
	}

	return NewMultiSigPublicKey(cryptography.MultiSigPublicKeyStruct{PubKeyMap: pubKeyMap, Threshold: threshold})
}

// Checks if two MultiSig public keys are equal
func (key *MultiSigPublicKey) Equals(publicKey cryptography.PublicKey) bool {
	return key.BasePublicKey.Equals(publicKey)
}

// Return the Sui address associated with this MultiSig public key
func (multisig *MultiSigPublicKey) ToSuiAddress() string {
	tmp := new(bytes.Buffer)
	tmp.WriteByte(cryptography.SignatureSchemeToFlag[cryptography.MultiSigScheme])

	threshold, err := bcs.Marshal(multisig.multisigPublicKey.Threshold)
	if err != nil {
		// TODO: throw an error?
		return ""
	}
	tmp.Write(threshold)

	for _, publicKey := range multisig.publicKeys {
		tmp.Write(publicKey.PublicKey.ToSuiBytes())
		tmp.WriteByte(publicKey.Weight)
	}

	sum256 := blake2b.Sum256(tmp.Bytes())
	return utils.NormalizeShortAddress(hex.EncodeToString(sum256[:])[:64])
}

// Return the byte array representation of the MultiSig public key
func (multisig *MultiSigPublicKey) ToRawBytes() []byte {
	return multisig.rawBytes
}

func (multisig *MultiSigPublicKey) GetPublicKeys() []cryptography.MultiSigPublicKeyPair {
	return multisig.publicKeys
}

func (multisig *MultiSigPublicKey) GetThreshold() uint16 {
	return multisig.multisigPublicKey.Threshold
}

// Return the Sui address associated with this MultiSig public key
func (multisig *MultiSigPublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.MultiSigScheme]
}

// Verifies that the signature is valid for for the provided message
func (multisig *MultiSigPublicKey) Verify(message []byte, multisigSignature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(multisigSignature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.MultiSigScheme {
		return false, err
	}

	thisMultisig := parsed.Multisig

	bs1, err := bcs.Marshal(multisig.multisigPublicKey)
	if err != nil {
		return false, err
	}
	bs2, err := bcs.Marshal(thisMultisig.MultisigPubKey)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(bs1, bs2) {
		return false, err
	}

	var signatureWeight uint16 = 0
	partialParsedData, err := ParsePartialSignatures(thisMultisig)
	if err != nil {
		return false, err
	}

	for _, data := range partialParsedData {
		signature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{
			SignatureScheme: cryptography.SignatureFlagToScheme[data.PublicKey.Flag()],
			Signature:       data.Signature,
			PublicKey:       data.PublicKey,
		})
		if err != nil {
			return false, err
		}

		pass, err := data.PublicKey.Verify(message, signature)
		if err != nil {
			return false, err
		}
		if !pass {
			return false, nil
		}

		signatureWeight += uint16(data.Weight)
	}
	return (signatureWeight >= multisig.multisigPublicKey.Threshold), nil
}

func (multisig *MultiSigPublicKey) CombinePartialSignatures(signatures []cryptography.SerializedSignature) (cryptography.SerializedSignature, error) {
	if len(signatures) > MaxSignerInMultisig {
		return "", fmt.Errorf("max number of signatures in a multisig is %d", MaxSignerInMultisig)
	}

	var bitmap uint16 = 0
	compressedSignatures := make([]cryptography.CompressedSignature, len(signatures))
	for i := 0; i < len(signatures); i++ {
		parsed, err := cryptography.ParseSerializedSignature(signatures[i])
		if err != nil {
			return "", err
		}

		if parsed.SignatureScheme == cryptography.MultiSigScheme {
			return "", fmt.Errorf("multisig is not supported inside MultiSig")
		}
		if parsed.SignatureScheme == cryptography.ZkLoginScheme {
			return "", fmt.Errorf("unimplemented %v", parsed.SignatureScheme)
		}

		tmp := new(bytes.Buffer)
		tmp.Write([]byte{cryptography.SignatureSchemeToFlag[parsed.SignatureScheme]})
		tmp.Write(parsed.Signature)
		compressedSignatures[i] = cryptography.CompressedSignature{Signature: [65]byte(tmp.Bytes())}

		var publicKeyIndex *int
		for j := 0; j < len(multisig.publicKeys); j++ {
			if bytes.Equal(parsed.PubKey, multisig.publicKeys[j].PublicKey.ToRawBytes()) {
				if bitmap&(1<<j) > 0 {
					return "", fmt.Errorf("received multiple signatures from the same public key")
				}
				publicKeyIndex = &j
				break
			}
		}

		if publicKeyIndex == nil {
			return "", fmt.Errorf("received signature from unknown public key")
		}
		bitmap |= 1 << *publicKeyIndex
	}

	m := cryptography.MultiSigStruct{
		Sigs:           compressedSignatures,
		Bitmap:         bitmap,
		MultisigPubKey: multisig.multisigPublicKey,
	}

	bs, err := bcs.Marshal(&m)
	if err != nil {
		return "", err
	}

	tmp := new(bytes.Buffer)
	tmp.Write([]byte{cryptography.SignatureSchemeToFlag[cryptography.MultiSigScheme]})
	tmp.Write(bs)

	return b64.ToBase64(tmp.Bytes()), nil
}

// Parse multisig structure into an array of individual signatures: signature scheme, the actual individual signature, public key and its weight.
func ParsePartialSignatures(multisig *cryptography.MultiSigStruct) ([]ParsedPartialMultiSigSignature, error) {
	res := make([]ParsedPartialMultiSigSignature, len(multisig.Sigs))

	for i := 0; i < len(multisig.Sigs); i++ {
		signatureScheme := cryptography.SignatureFlagToScheme[multisig.Sigs[i].Signature[0]]
		signature := multisig.Sigs[i].Signature[1:]

		bitmapIndices, err := AsIndices(multisig.Bitmap)
		if err != nil {
			return nil, err
		}
		pkIndex := bitmapIndices[i]
		pair := multisig.MultisigPubKey.PubKeyMap[pkIndex]
		pkBytes := pair.PubKey[1:]

		if signatureScheme == cryptography.MultiSigScheme {
			return nil, fmt.Errorf("multisig is not supported inside MultiSig")
		}

		publicKey, err := verify.PublicKeyFromRawBytes(signatureScheme, pkBytes)
		if err != nil {
			return nil, err
		}

		res[i] = ParsedPartialMultiSigSignature{
			SignatureScheme: signatureScheme,
			Signature:       signature,
			PublicKey:       publicKey,
			Weight:          pair.Weight,
		}

	}
	return res, nil
}

func AsIndices(bitmap uint16) ([]byte, error) {
	if bitmap > 1024 {
		return nil, fmt.Errorf("invalid bitmap")
	}

	res := []byte{}
	for i := 0; i < 10; i++ {
		if (bitmap & (1 << i)) != 0 {
			res = append(res, byte(i))
		}
	}
	return res, nil
}

func PublicKeyFromSuiBytes[T string | []byte](publicKey T) (pk cryptography.PublicKey, err error) {
	var bs []byte
	switch v := any(publicKey).(type) {
	case string:
		bs, err = b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		bs = v
	}

	signatureScheme := cryptography.SignatureFlagToScheme[bs[0]]

	if signatureScheme == cryptography.ZkLoginScheme {
		return nil, fmt.Errorf("zkLogin publicKey is not supported")
	}
	return verify.PublicKeyFromRawBytes(signatureScheme, bs[1:])
}
