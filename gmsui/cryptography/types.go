package cryptography

type CompressedSignature struct {
	Signature [65]byte `json:"signature"`
}

type PubKeyEnumWeightPair struct {
	PubKey [33]byte `json:"pubKey"`
	Weight uint8    `json:"weight"`
}

type StringPubKeyEnumWeightPair struct {
	PubKey string `json:"pubKey"`
	Weight uint8  `json:"weight"`
}

type MultiSigPublicKeyStruct struct {
	PubKeyMap []PubKeyEnumWeightPair `json:"pubKeymap"`
	Threshold uint16                 `json:"threshold"`
}

type MultiSigStruct struct {
	Sigs           []CompressedSignature   `json:"sigs"`
	Bitmap         uint16                  `json:"bitmap"`
	MultisigPubKey MultiSigPublicKeyStruct `json:"multisigPubKey"`
}
