package gmsui

type SuiSignedTransactionRet = SuiSignedDataRet
type SuiSignedMessageRet = SuiSignedDataRet

type SuiSignedDataRet struct {
	TxBytes   string `json:"tx_bytes"`
	Signature string `json:"signature"`
}

type SignaturePubkeyPair struct {
	SignatureScheme string
	Signature       []byte
	PubKey          []byte
}
