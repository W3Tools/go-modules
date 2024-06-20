package cryptography

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/tyler-smith/go-bip39"
)

// Parse and validate a path that is compliant to SLIP-0010 in form m/44'/784'/{account_index}'/{change_index}'/{address_index}'.
// e.g. `m/44'/784'/0'/0'/0'`
func IsValidHardenedPath(path string) bool {
	matched, _ := regexp.MatchString("^m\\/44'\\/784'\\/[0-9]+'\\/[0-9]+'\\/[0-9]+'+$", path)
	return matched
}

// Parse and validate a path that is compliant to BIP-32 in form m/54'/784'/{account_index}'/{change_index}/{address_index}
// for Secp256k1 and m/74'/784'/{account_index}'/{change_index}/{address_index} for Secp256r1.
// Note that the purpose for Secp256k1 is registered as 54, to differentiate from Ed25519 with purpose 44.
// e.g. `m/54'/784'/0'/0/0`
func IsValidBIP32Path(path string) bool {
	matched, _ := regexp.MatchString("^m\\/(54|74)'\\/784'\\/[0-9]+'\\/[0-9]+\\/[0-9]+$", path)
	return matched
}

// Uses KDF to derive 64 bytes of key data from mnemonic with empty password.
// mnemonics 12 words string split by spaces.
func MnemonicToSeed(mnemonics string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonics) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	return bip39.NewSeed(mnemonics, ""), nil
}

// Derive the seed in hex format from a 12-word mnemonic string.
// mnemonics 12 words string split by spaces.
func MnemonicToSeedHex(mnemonics string) (string, error) {
	seed, err := MnemonicToSeed(mnemonics)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(seed), nil
}

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}
