package multisig

import (
	"fmt"
	"strings"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/W3Tools/go-modules/gmsui/keypairs/ed25519"
)

var (
	keypair0          *ed25519.Ed25519PublicKey
	keypair1          *ed25519.Ed25519PublicKey
	keypair2          *ed25519.Ed25519PublicKey
	keypair3          *ed25519.Ed25519PublicKey
	multisigPublicKey *MultiSigPublicKey
)

func init() {
	var err error
	keypair0, err = ed25519.NewEd25519PublicKey("RIIlNygk0BhV4AyP1Ntwa9jeiz9MULhjODIEvWSJ0gk=")
	if err != nil {
		panic(fmt.Sprintf("new ed25519 public key err, msg: %v", err))
	}
	fmt.Printf(`keypair0 -> address: %v
	    flag: %v
	    toRawBytes: %v
	    toBase64: %v
	    toSuiBytes: %v
	    toSuiPublicKey: %v`,
		keypair0.ToSuiAddress(), keypair0.Flag(), keypair0.ToRawBytes(), keypair0.ToBase64(), keypair0.ToSuiBytes(), keypair0.ToSuiPublicKey())
	fmt.Println("")

	keypair1, err = ed25519.NewEd25519PublicKey("H8gKV7gmH1T7YmAfwp5FNPWas0K0AFY2OBQ+0uOY57E=")
	if err != nil {
		panic(fmt.Sprintf("new ed25519 public key err, msg: %v", err))
	}
	fmt.Printf(`keypair1 -> address: %v
	    flag: %v
	    toRawBytes: %v
	    toBase64: %v
	    toSuiBytes: %v
	    toSuiPublicKey: %v`,
		keypair1.ToSuiAddress(), keypair1.Flag(), keypair1.ToRawBytes(), keypair1.ToBase64(), keypair1.ToSuiBytes(), keypair1.ToSuiPublicKey())
	fmt.Println("")

	keypair2, err = ed25519.NewEd25519PublicKey("hGxdp+5Molt7N1VRSvp5eExYViexZzISsYW51okg+7w=")
	if err != nil {
		panic(fmt.Sprintf("new ed25519 public key err, msg: %v", err))
	}
	fmt.Printf(`keypair2 -> address: %v
	    flag: %v
	    toRawBytes: %v
	    toBase64: %v
	    toSuiBytes: %v
	    toSuiPublicKey: %v`,
		keypair2.ToSuiAddress(), keypair2.Flag(), keypair2.ToRawBytes(), keypair2.ToBase64(), keypair2.ToSuiBytes(), keypair2.ToSuiPublicKey())
	fmt.Println("")

	keypair3, err = ed25519.NewEd25519PublicKey("bkIKm0xtKG3cPNZsFnIy+HDmBp8kLQspBu1rWqbmRbE=")
	if err != nil {
		panic(fmt.Sprintf("new ed25519 public key err, msg: %v", err))
	}
	fmt.Printf(`keypair3 -> address: %v
	    flag: %v
	    toRawBytes: %v
	    toBase64: %v
	    toSuiBytes: %v
	    toSuiPublicKey: %v`,
		keypair3.ToSuiAddress(), keypair3.Flag(), keypair3.ToRawBytes(), keypair3.ToBase64(), keypair3.ToSuiBytes(), keypair3.ToSuiPublicKey())
	fmt.Println("")

	multisigPublicKey, err = new(MultiSigPublicKey).FromPublicKeys(
		[]PublicKeyWeightPair{
			{
				PublicKey: keypair0,
				Weight:    1,
			},
			{
				PublicKey: keypair1,
				Weight:    1,
			},
			{
				PublicKey: keypair2,
				Weight:    1,
			},
			{
				PublicKey: keypair3,
				Weight:    1,
			},
		},
		3,
	)
	if err != nil {
		panic(fmt.Sprintf("new multisig public key from public keys error, msg: %v", err))
	}

	fmt.Printf(`multisig -> address: %v
	    flag: %v
	    threshold: %v
	    base64: %v
	    raw bytes: %v
	    sui bytes: %v
	    sui public key: %v`,
		multisigPublicKey.ToSuiAddress(), multisigPublicKey.Flag(), multisigPublicKey.GetThreshold(), multisigPublicKey.ToBase64(), multisigPublicKey.ToRawBytes(), multisigPublicKey.ToSuiBytes(), multisigPublicKey.ToSuiPublicKey())
	fmt.Println("")
	fmt.Printf("%s\n", strings.Repeat("-", 150))
}

func TestCombinePartialSignatures(t *testing.T) {
	signature0 := "AJe72Zj/9Ooyv8mHpZjzKXEF6NLbeQKr/pfWDPfCi7yofyY6s0WRb/z+/VqEfg/GhRpmI7dHEARGjXBOxGBTSghEgiU3KCTQGFXgDI/U23Br2N6LP0xQuGM4MgS9ZInSCQ=="
	signature2 := "AKdszUyALLWYUz7XK9Hm6FYJ0d3quIZ39wSx1GUUCCI4HNuNbTer+3v0gf6+dFaWs+6JySl87OJbZjYVmJaqWAuEbF2n7kyiW3s3VVFK+nl4TFhWJ7FnMhKxhbnWiSD7vA=="
	signature3 := "AN/ugnW8oV1qLZptKuGx7G+yH2bhtTeK7OPe1wj0SHvYp74T/3ccWhaA0+Qzy0k5eC1EeRkndGgjGp7XA0LZrQxuQgqbTG0obdw81mwWcjL4cOYGnyQtCykG7WtapuZFsQ=="

	data, err := multisigPublicKey.CombinePartialSignatures([]cryptography.SerializedSignature{signature0, signature3, signature2})
	if err != nil {
		t.Fatalf("combine partial signature err, msg: %v", err)
	}
	fmt.Printf("combinePartialSignatures: %v\n", data)

	txb := "AAABACAsFID88X4wXT/z/MbNiECsiW73FJO3HyyKXhySS1OBZgIEAdwCoRzrCwYAAAAKAQAIAggMAxQVBCkCBSsVB0BdCJ0BQArdAQgM5QFJDa4CAgADAQgBCwEMAAAIAAECBAADAQIAAAUAAQAACgIBAAEGAAMAAgkFAQEIAwQBBwgCAAIHCAADAQgBAQgAAQkAAQMHQ291bnRlcglUeENvbnRleHQDVUlEBWhlbGxvAmlkBGluaXQDbmV3Bm51bWJlcgZvYmplY3QMc2hhcmVfb2JqZWN0BXRvdWNoCHRyYW5zZmVyCnR4X2NvbnRleHQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAICBAgBBwMAAAAAAQYLABECBgAAAAAAAAAAEgA4AAIBAQAABhEKAQwCCwEGAAAAAAAAAAAhBAgGAQAAAAAAAAAMAgoAEAAUCwIWCwAPABUCAAEAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIBAQIAAAEAACwUgPzxfjBdP/P8xs2IQKyJbvcUk7cfLIpeHJJLU4FmAbLXqgfMag7gldu5IkqvV+goaqyaVbeXAVPw2SgCGVoVUJ42DwAAAAAg2JaDc2hq9u0qw503zPRRv39t6ksRcSwKKV8Bom3GtwIsFID88X4wXT/z/MbNiECsiW73FJO3HyyKXhySS1OBZu4CAAAAAAAAgJaYAAAAAAAA"
	b64txb, err := b64.FromBase64(txb)
	if err != nil {
		t.Fatalf("from base64 err, msg: %v", err)
	}

	pass, err := multisigPublicKey.VerifyTransactionBlock(b64txb, data)
	if err != nil {
		t.Fatalf("verify transaction block err, msg: %v", err)
	}
	fmt.Printf("pass: %v\n", pass)
}

func TestVerifyPersonalMessage(t *testing.T) {
	message := []byte("hello world")

	signature1 := "AH89ba1KO7KoMQfjGJOkxKaI80GseFWQ/GKTF7pV0hc4LpAnHZ0lfCQWavvTL5v6DaIzTBYlpD9V7kFesCr80whEgiU3KCTQGFXgDI/U23Br2N6LP0xQuGM4MgS9ZInSCQ=="
	data, err := multisigPublicKey.CombinePartialSignatures([]cryptography.SerializedSignature{signature1})
	if err != nil {
		t.Fatalf("combine partial signature err, msg: %v", err)
	}
	fmt.Printf("combine partial signatures: %v\n", data)

	pass, err := multisigPublicKey.VerifyPersonalMessage(message, data)
	if err != nil {
		t.Fatalf("verify personal message err, msg: %v", err)
	}
	fmt.Printf("pass: %v\n", pass)
}

func TestVerifyTransactionBlock(t *testing.T) {
	signature := "AH89ba1KO7KoMQfjGJOkxKaI80GseFWQ/GKTF7pV0hc4LpAnHZ0lfCQWavvTL5v6DaIzTBYlpD9V7kFesCr80whEgiU3KCTQGFXgDI/U23Br2N6LP0xQuGM4MgS9ZInSCQ=="
	txb := "AAABACAsFID88X4wXT/z/MbNiECsiW73FJO3HyyKXhySS1OBZgIEAdwCoRzrCwYAAAAKAQAIAggMAxQVBCkCBSsVB0BdCJ0BQArdAQgM5QFJDa4CAgADAQgBCwEMAAAIAAECBAADAQIAAAUAAQAACgIBAAEGAAMAAgkFAQEIAwQBBwgCAAIHCAADAQgBAQgAAQkAAQMHQ291bnRlcglUeENvbnRleHQDVUlEBWhlbGxvAmlkBGluaXQDbmV3Bm51bWJlcgZvYmplY3QMc2hhcmVfb2JqZWN0BXRvdWNoCHRyYW5zZmVyCnR4X2NvbnRleHQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAICBAgBBwMAAAAAAQYLABECBgAAAAAAAAAAEgA4AAIBAQAABhEKAQwCCwEGAAAAAAAAAAAhBAgGAQAAAAAAAAAMAgoAEAAUCwIWCwAPABUCAAEAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIBAQIAAAEAACwUgPzxfjBdP/P8xs2IQKyJbvcUk7cfLIpeHJJLU4FmAbLXqgfMag7gldu5IkqvV+goaqyaVbeXAVPw2SgCGVoVUJ42DwAAAAAg2JaDc2hq9u0qw503zPRRv39t6ksRcSwKKV8Bom3GtwIsFID88X4wXT/z/MbNiECsiW73FJO3HyyKXhySS1OBZu4CAAAAAAAAgJaYAAAAAAAA"

	b64txb, err := b64.FromBase64(txb)
	if err != nil {
		t.Fatalf("from base64 err, msg: %v", err)
	}
	pass, err := keypair0.VerifyTransactionBlock(b64txb, signature)
	if err != nil {
		t.Fatalf("verify transaction block err, msg: %v", err)
	}
	fmt.Printf("pass: %v\n", pass)
}
