package cryptography

type SignatureScheme = string
type SignatureFlag = uint8

var (
	Ed25519Scheme   SignatureScheme = "ED25519"
	Secp256k1Scheme SignatureScheme = "Secp256k1"
	Secp256r1Scheme SignatureScheme = "Secp256r1"
	MultiSigScheme  SignatureScheme = "MultiSig"
	ZkLoginScheme   SignatureScheme = "ZkLogin"
)

var SignatureSchemeToFlag = map[SignatureScheme]SignatureFlag{
	Ed25519Scheme:   0x00,
	Secp256k1Scheme: 0x01,
	Secp256r1Scheme: 0x02,
	MultiSigScheme:  0x03,
	ZkLoginScheme:   0x05,
}

var SignatureSchemeToSize = map[SignatureScheme]int{
	Ed25519Scheme:   32,
	Secp256k1Scheme: 33,
	Secp256r1Scheme: 33,
}

var SignatureFlagToScheme = map[SignatureFlag]SignatureScheme{
	0x00: Ed25519Scheme,
	0x01: Secp256k1Scheme,
	0x02: Secp256r1Scheme,
	0x03: MultiSigScheme,
	0x05: ZkLoginScheme,
}
