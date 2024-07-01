package gmsui

type SignaturePubkeyPair struct {
	SignatureScheme string
	Signature       []byte
	PubKey          []byte
}

type MultiSigStruct struct {
	Sigs       []CompressedSignature   `json:"sigs"`
	Bitmap     uint16                  `json:"bitmap"`
	MultisigPK MultiSigPublicKeyStruct `json:"multisig_pk"`
}

type CompressedSignature struct {
	Signature [65]byte `json:"signature"`
}

type MultiSigPublicKeyStruct struct {
	PKMap     []PubkeyEnumWeightPair `json:"pk_map"`
	Threshold uint16                 `json:"threshold"`
}

type PubkeyEnumWeightPair struct {
	PubKey [33]byte `json:"pub_key"`
	Weight uint8    `json:"weight"`
}

type SuiMultiSigInfo struct {
	Address   string
	Threshold uint16
	Signers   []SuiMultiSigInfoSigner
}

type SuiMultiSigInfoSigner struct {
	Address      string
	B64PublicKey string
	HexPublicKey string
	Weight       uint8
}
