package tests

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/W3Tools/go-modules/gmsui/keypairs/ed25519"
	"github.com/W3Tools/go-modules/gmsui/keypairs/secp256k1"
	"github.com/W3Tools/go-modules/gmsui/keypairs/secp256r1"
	"github.com/W3Tools/go-modules/gmsui/multisig"
)

const mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
const message = "Hello Go Modules!"

func TestEd25519(t *testing.T) {
	testDatas := []struct {
		path     string
		expected struct {
			getKeyScheme            string
			getSecretKey            string
			toSuiAddress            string
			publicKeyToRawBytes     []byte
			publicKeyToSuiBytes     []byte
			publicKeyToBase64       string
			publicKeyToSuiPublicKey string
			publicKeyToSuiAddress   string
			message                 string
			signature               string
		}
	}{
		{
			path: "m/44'/784'/0'/0'/0'",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "ED25519",
				getSecretKey:            "suiprivkey1qzyxnjc8z79lvlsg6lz2hh69fp7m7duunfzjlnkzsd59f062855mqacydfr",
				toSuiAddress:            "0x5e93a736d04fbb25737aa40bee40171ef79f65fae833749e3c089fe7cc2161f1",
				publicKeyToRawBytes:     []byte{144, 11, 77, 129, 238, 206, 163, 223, 47, 116, 177, 66, 0, 196, 244, 207, 63, 73, 175, 172, 167, 166, 52, 255, 210, 207, 111, 248, 43, 218, 236, 242},
				publicKeyToSuiBytes:     []byte{0, 144, 11, 77, 129, 238, 206, 163, 223, 47, 116, 177, 66, 0, 196, 244, 207, 63, 73, 175, 172, 167, 166, 52, 255, 210, 207, 111, 248, 43, 218, 236, 242},
				publicKeyToBase64:       "kAtNge7Oo98vdLFCAMT0zz9Jr6ynpjT/0s9v+Cva7PI=",
				publicKeyToSuiPublicKey: "AJALTYHuzqPfL3SxQgDE9M8/Sa+sp6Y0/9LPb/gr2uzy",
				publicKeyToSuiAddress:   "0x5e93a736d04fbb25737aa40bee40171ef79f65fae833749e3c089fe7cc2161f1",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AGtXLcPTNs1EukLef73WVQ+Q0P+9uyrbu/g4u3X4H/uCgbhk3G6Dg46xO9Bs5C78wcmqE9p1sZO0UWsg0l5UrwGQC02B7s6j3y90sUIAxPTPP0mvrKemNP/Sz2/4K9rs8g==",
			},
		},
		{
			path: "m/44'/784'/0'/0'/1'",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "ED25519",
				getSecretKey:            "suiprivkey1qpqk4vck5utyfttht7yphs4xc5flc7ymvzts36vh077f27d9lj942dj2rlk",
				toSuiAddress:            "0xf7c7a39996ac7f1c307b96c96d65cce0855dcc7ccd021c453964f2f62f98e71f",
				publicKeyToRawBytes:     []byte{72, 13, 240, 13, 190, 79, 51, 38, 217, 189, 169, 144, 121, 149, 79, 166, 157, 83, 121, 195, 78, 102, 210, 173, 163, 102, 215, 119, 167, 200, 112, 93},
				publicKeyToSuiBytes:     []byte{0, 72, 13, 240, 13, 190, 79, 51, 38, 217, 189, 169, 144, 121, 149, 79, 166, 157, 83, 121, 195, 78, 102, 210, 173, 163, 102, 215, 119, 167, 200, 112, 93},
				publicKeyToBase64:       "SA3wDb5PMybZvamQeZVPpp1TecNOZtKto2bXd6fIcF0=",
				publicKeyToSuiPublicKey: "AEgN8A2+TzMm2b2pkHmVT6adU3nDTmbSraNm13enyHBd",
				publicKeyToSuiAddress:   "0xf7c7a39996ac7f1c307b96c96d65cce0855dcc7ccd021c453964f2f62f98e71f",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AOwzsOUKlYyE9140S59Gw/giW6AWRTGDH2qhCxoBXa13cBlLyUP2y+4mh2MTGZbl8jdE4dxQmB+fez9UqIFXdAFIDfANvk8zJtm9qZB5lU+mnVN5w05m0q2jZtd3p8hwXQ==",
			},
		},
		{
			path: "m/44'/784'/0'/0'/100'",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "ED25519",
				getSecretKey:            "suiprivkey1qrf7g3kgaa0e9qrdk5hkgn9hgu8ns35eqynxrmv39dltsf07ksrny7hmmm7",
				toSuiAddress:            "0x09bc557f22f2a7d19dbbb2e0862164e8f119d1a085356458e25679d2ece2fbe7",
				publicKeyToRawBytes:     []byte{181, 33, 229, 132, 252, 227, 116, 97, 174, 182, 8, 43, 122, 79, 119, 187, 164, 229, 102, 2, 163, 232, 176, 67, 77, 126, 236, 43, 254, 24, 159, 92},
				publicKeyToSuiBytes:     []byte{0, 181, 33, 229, 132, 252, 227, 116, 97, 174, 182, 8, 43, 122, 79, 119, 187, 164, 229, 102, 2, 163, 232, 176, 67, 77, 126, 236, 43, 254, 24, 159, 92},
				publicKeyToBase64:       "tSHlhPzjdGGutggrek93u6TlZgKj6LBDTX7sK/4Yn1w=",
				publicKeyToSuiPublicKey: "ALUh5YT843RhrrYIK3pPd7uk5WYCo+iwQ01+7Cv+GJ9c",
				publicKeyToSuiAddress:   "0x09bc557f22f2a7d19dbbb2e0862164e8f119d1a085356458e25679d2ece2fbe7",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AGpkrflwgasF/JEju0E+alpK0lIYw7a3JQj1HCM121PsF8W6yRcowdSQo+z7cwnQ7dhZrmDRvvAI30EICYI2PQm1IeWE/ON0Ya62CCt6T3e7pOVmAqPosENNfuwr/hifXA==",
			},
		},
	}
	for _, test := range testDatas {
		t.Run(test.path, func(t *testing.T) {
			// For keypair
			keypair, err := ed25519.DeriveKeypair(mnemonic, test.path)
			if err != nil {
				t.Fatalf("failed to derive ed25519 keypair, msg: %v", err)
			}

			getKeyScheme := keypair.GetKeyScheme()
			if !reflect.DeepEqual(getKeyScheme, test.expected.getKeyScheme) {
				t.Errorf("unable to get key scheme, expected %s, but got %s", test.expected.getKeyScheme, getKeyScheme)
			}

			getSecretKey, err := keypair.GetSecretKey()
			if err != nil {
				t.Fatalf("failed to get secret key, msg: %v", err)
			}
			if !reflect.DeepEqual(getSecretKey, test.expected.getSecretKey) {
				t.Errorf("unable to get secret key, expected %v, but got %v", test.expected.getSecretKey, getSecretKey)
			}

			if !reflect.DeepEqual(keypair.ToSuiAddress(), test.expected.toSuiAddress) {
				t.Errorf("unable to sui address, expected %v, but got %v", test.expected.toSuiAddress, keypair.ToSuiAddress())
			}

			// For public key
			pubkey, err := keypair.GetPublicKey()
			if err != nil {
				t.Fatalf("failed to get public key, msg: %v", err)
			}
			if !bytes.Equal(pubkey.ToRawBytes(), test.expected.publicKeyToRawBytes) {
				t.Errorf("unable to get public key to raw bytes, expected %v, but got %v", test.expected.publicKeyToRawBytes, pubkey.ToRawBytes())
			}

			if !bytes.Equal(pubkey.ToSuiBytes(), test.expected.publicKeyToSuiBytes) {
				t.Errorf("unable to get public key to sui bytes, expected %v, but got %v", test.expected.publicKeyToSuiBytes, pubkey.ToSuiBytes())
			}

			if !reflect.DeepEqual(pubkey.ToBase64(), test.expected.publicKeyToBase64) {
				t.Errorf("unable to get public key to base64, expected %v, but got %v", test.expected.publicKeyToBase64, pubkey.ToBase64())
			}

			if !reflect.DeepEqual(pubkey.ToSuiPublicKey(), test.expected.publicKeyToSuiPublicKey) {
				t.Errorf("unable to get public key to sui public key, expected %v, but got %v", test.expected.publicKeyToSuiPublicKey, pubkey.ToSuiPublicKey())
			}

			if !reflect.DeepEqual(pubkey.ToSuiAddress(), test.expected.publicKeyToSuiAddress) {
				t.Errorf("unable to get public key to sui address, expected %v, but got %v", test.expected.publicKeyToSuiAddress, pubkey.ToSuiAddress())
			}

			// For signature
			data, err := keypair.SignPersonalMessage([]byte(message))
			if err != nil {
				t.Fatalf("keypair failed to sign personal message, msg: %v", err)
			}

			if !reflect.DeepEqual(data.Bytes, test.expected.message) {
				t.Errorf("unable to sign personal message, bytes expected %s, but got %v", test.expected.message, data.Bytes)
			}

			if !reflect.DeepEqual(data.Signature, test.expected.signature) {
				t.Errorf("unable to sign personal message, signature expected %s, but got %v", test.expected.signature, data.Signature)
			}

			pass, err := pubkey.VerifyPersonalMessage([]byte(message), data.Signature)
			if err != nil {
				t.Fatalf("public key failed to verify personal message, msg: %v", err)
			}
			if !pass {
				t.Errorf("unable to verify personal message, expected %v, but got %v", true, pass)
			}
		})
	}
}

func TestSecp256k1(t *testing.T) {
	testDatas := []struct {
		path     string
		expected struct {
			getKeyScheme            string
			getSecretKey            string
			toSuiAddress            string
			publicKeyToRawBytes     []byte
			publicKeyToSuiBytes     []byte
			publicKeyToBase64       string
			publicKeyToSuiPublicKey string
			publicKeyToSuiAddress   string
			message                 string
			signature               string
		}
	}{
		{
			path: "m/54'/784'/0'/0/0",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "Secp256k1",
				getSecretKey:            "suiprivkey1qy82eu8yuzp4dykhe5d8cth23sw05yxnqqzpf5ce0rnmdjn905rgge6hyhd",
				toSuiAddress:            "0xc61a7f1161020a717f852dca2e9bfc1ffe235145406dfbdccc16e6907c1f5403",
				publicKeyToRawBytes:     []byte{2, 98, 61, 134, 15, 70, 204, 233, 17, 125, 63, 26, 195, 130, 183, 156, 89, 146, 138, 0, 74, 25, 134, 86, 26, 153, 223, 42, 133, 22, 124, 245, 133},
				publicKeyToSuiBytes:     []byte{1, 2, 98, 61, 134, 15, 70, 204, 233, 17, 125, 63, 26, 195, 130, 183, 156, 89, 146, 138, 0, 74, 25, 134, 86, 26, 153, 223, 42, 133, 22, 124, 245, 133},
				publicKeyToBase64:       "AmI9hg9GzOkRfT8aw4K3nFmSigBKGYZWGpnfKoUWfPWF",
				publicKeyToSuiPublicKey: "AQJiPYYPRszpEX0/GsOCt5xZkooAShmGVhqZ3yqFFnz1hQ==",
				publicKeyToSuiAddress:   "0xc61a7f1161020a717f852dca2e9bfc1ffe235145406dfbdccc16e6907c1f5403",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AbtKlpY/Bsmo9huj2TiGdD92phTWxx3ABn4t/McFV7iGQFwzhvj8loW95rvoXplGC5XvrERwLk9XPYNpS9K758sCYj2GD0bM6RF9PxrDgrecWZKKAEoZhlYamd8qhRZ89YU=",
			},
		},
		{
			path: "m/54'/784'/0'/0/1",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "Secp256k1",
				getSecretKey:            "suiprivkey1qy7nhapqvq4jw3y9a2nr0sr8ug5awejjgs0d9zeq39hzh9afxr2u757dqhd",
				toSuiAddress:            "0x03de9efda2d82b61535b6f8448ea1ef55f914994f4b27f4628f918a054e55ba4",
				publicKeyToRawBytes:     []byte{2, 56, 161, 184, 104, 161, 161, 222, 177, 157, 85, 123, 132, 0, 169, 250, 46, 20, 141, 54, 137, 124, 85, 2, 113, 226, 87, 216, 253, 178, 5, 141, 81},
				publicKeyToSuiBytes:     []byte{1, 2, 56, 161, 184, 104, 161, 161, 222, 177, 157, 85, 123, 132, 0, 169, 250, 46, 20, 141, 54, 137, 124, 85, 2, 113, 226, 87, 216, 253, 178, 5, 141, 81},
				publicKeyToBase64:       "AjihuGihod6xnVV7hACp+i4UjTaJfFUCceJX2P2yBY1R",
				publicKeyToSuiPublicKey: "AQI4obhooaHesZ1Ve4QAqfouFI02iXxVAnHiV9j9sgWNUQ==",
				publicKeyToSuiAddress:   "0x03de9efda2d82b61535b6f8448ea1ef55f914994f4b27f4628f918a054e55ba4",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AUnwfiejszYSZ/2vP8+YrkcsP18tmNbV6Crqg1yV9YjDTvhBhVRTkupJvblaESJXgyWWBKEnZY4avNMJ/ZgwtTsCOKG4aKGh3rGdVXuEAKn6LhSNNol8VQJx4lfY/bIFjVE=",
			},
		},
		{
			path: "m/54'/784'/0'/0/100",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "Secp256k1",
				getSecretKey:            "suiprivkey1q9kmy20jqv70xh3arz7zq6e8xxcqncakde2zf4axc43a2wpw9f59cf0gjvh",
				toSuiAddress:            "0x968c9a3409ec574e66fa2275a41b349d33f286c1f3f9dab05cfe1dc0385ce56e",
				publicKeyToRawBytes:     []byte{2, 41, 208, 210, 198, 3, 116, 32, 177, 207, 164, 84, 229, 197, 122, 6, 197, 83, 107, 42, 185, 24, 176, 145, 21, 149, 60, 32, 175, 172, 196, 114, 33},
				publicKeyToSuiBytes:     []byte{1, 2, 41, 208, 210, 198, 3, 116, 32, 177, 207, 164, 84, 229, 197, 122, 6, 197, 83, 107, 42, 185, 24, 176, 145, 21, 149, 60, 32, 175, 172, 196, 114, 33},
				publicKeyToBase64:       "AinQ0sYDdCCxz6RU5cV6BsVTayq5GLCRFZU8IK+sxHIh",
				publicKeyToSuiPublicKey: "AQIp0NLGA3Qgsc+kVOXFegbFU2squRiwkRWVPCCvrMRyIQ==",
				publicKeyToSuiAddress:   "0x968c9a3409ec574e66fa2275a41b349d33f286c1f3f9dab05cfe1dc0385ce56e",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "ASZEpic8S267G1WS6Tu/0eNdeCsyjOTmD6DVImGP19w2dij/Px6FKrL4tGHtQlUJbf/zVCJvJVYD0QEtR9d24LICKdDSxgN0ILHPpFTlxXoGxVNrKrkYsJEVlTwgr6zEciE=",
			},
		},
	}
	for _, test := range testDatas {
		t.Run(test.path, func(t *testing.T) {
			// For keypair
			keypair, err := secp256k1.DeriveKeypair(mnemonic, test.path)
			if err != nil {
				t.Fatalf("failed to derive secp256k1 keypair, msg: %v", err)
			}

			getKeyScheme := keypair.GetKeyScheme()
			if !reflect.DeepEqual(getKeyScheme, test.expected.getKeyScheme) {
				t.Errorf("unable to get key scheme, expected %s, but got %s", test.expected.getKeyScheme, getKeyScheme)
			}

			getSecretKey, err := keypair.GetSecretKey()
			if err != nil {
				t.Fatalf("failed to get secret key, msg: %v", err)
			}
			if !reflect.DeepEqual(getSecretKey, test.expected.getSecretKey) {
				t.Errorf("unable to get secret key, expected %v, but got %v", test.expected.getSecretKey, getSecretKey)
			}

			if !reflect.DeepEqual(keypair.ToSuiAddress(), test.expected.toSuiAddress) {
				t.Errorf("unable to sui address, expected %v, but got %v", test.expected.toSuiAddress, keypair.ToSuiAddress())
			}

			// For public key
			pubkey, err := keypair.GetPublicKey()
			if err != nil {
				t.Fatalf("failed to get public key, msg: %v", err)
			}
			if !bytes.Equal(pubkey.ToRawBytes(), test.expected.publicKeyToRawBytes) {
				t.Errorf("unable to get public key to raw bytes, expected %v, but got %v", test.expected.publicKeyToRawBytes, pubkey.ToRawBytes())
			}

			if !bytes.Equal(pubkey.ToSuiBytes(), test.expected.publicKeyToSuiBytes) {
				t.Errorf("unable to get public key to sui bytes, expected %v, but got %v", test.expected.publicKeyToSuiBytes, pubkey.ToSuiBytes())
			}

			if !reflect.DeepEqual(pubkey.ToBase64(), test.expected.publicKeyToBase64) {
				t.Errorf("unable to get public key to base64, expected %v, but got %v", test.expected.publicKeyToBase64, pubkey.ToBase64())
			}

			if !reflect.DeepEqual(pubkey.ToSuiPublicKey(), test.expected.publicKeyToSuiPublicKey) {
				t.Errorf("unable to get public key to sui public key, expected %v, but got %v", test.expected.publicKeyToSuiPublicKey, pubkey.ToSuiPublicKey())
			}

			if !reflect.DeepEqual(pubkey.ToSuiAddress(), test.expected.publicKeyToSuiAddress) {
				t.Errorf("unable to get public key to sui address, expected %v, but got %v", test.expected.publicKeyToSuiAddress, pubkey.ToSuiAddress())
			}

			// For signature
			data, err := keypair.SignPersonalMessage([]byte(message))
			if err != nil {
				t.Fatalf("keypair failed to sign personal message, msg: %v", err)
			}

			if !reflect.DeepEqual(data.Bytes, test.expected.message) {
				t.Errorf("unable to sign personal message, bytes expected %s, but got %v", test.expected.message, data.Bytes)
			}

			if !reflect.DeepEqual(data.Signature, test.expected.signature) {
				t.Errorf("unable to sign personal message, signature expected %s, but got %v", test.expected.signature, data.Signature)
			}

			pass, err := pubkey.VerifyPersonalMessage([]byte(message), data.Signature)
			if err != nil {
				t.Fatalf("public key failed to verify personal message, msg: %v", err)
			}
			if !pass {
				t.Errorf("unable to verify personal message, expected %v, but got %v", true, pass)
			}
		})
	}
}

func TestSecp256r1(t *testing.T) {
	testDatas := []struct {
		path     string
		expected struct {
			getKeyScheme            string
			getSecretKey            string
			toSuiAddress            string
			publicKeyToRawBytes     []byte
			publicKeyToSuiBytes     []byte
			publicKeyToBase64       string
			publicKeyToSuiPublicKey string
			publicKeyToSuiAddress   string
			message                 string
			signature               string
		}
	}{
		{
			path: "m/74'/784'/0'/0/0",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "Secp256r1",
				getSecretKey:            "suiprivkey1qfy8w9uvgleu04l6um3u6dagtsdgy8mysx8rr4qlyjn9ummcg34dzjn2gkr",
				toSuiAddress:            "0x0c0f9f53f2ad697e18279dfadefdd070c8e99416309d3ce614086c0860db6bb4",
				publicKeyToRawBytes:     []byte{3, 64, 25, 188, 168, 168, 120, 69, 138, 99, 229, 191, 83, 243, 8, 85, 227, 16, 112, 247, 181, 124, 249, 220, 242, 101, 201, 139, 219, 23, 187, 23, 196},
				publicKeyToSuiBytes:     []byte{2, 3, 64, 25, 188, 168, 168, 120, 69, 138, 99, 229, 191, 83, 243, 8, 85, 227, 16, 112, 247, 181, 124, 249, 220, 242, 101, 201, 139, 219, 23, 187, 23, 196},
				publicKeyToBase64:       "A0AZvKioeEWKY+W/U/MIVeMQcPe1fPnc8mXJi9sXuxfE",
				publicKeyToSuiPublicKey: "AgNAGbyoqHhFimPlv1PzCFXjEHD3tXz53PJlyYvbF7sXxA==",
				publicKeyToSuiAddress:   "0x0c0f9f53f2ad697e18279dfadefdd070c8e99416309d3ce614086c0860db6bb4",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AkgYRN9hEX5LgSlT+r/M/15e9UKJmGxFeUc+q4ozTgXzCUOkBXHdVGrqKrTm4M50wp/pAgNnnASSJVRnGSmjA14DQBm8qKh4RYpj5b9T8whV4xBw97V8+dzyZcmL2xe7F8Q=",
			},
		},
		{
			path: "m/74'/784'/0'/0/1",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "Secp256r1",
				getSecretKey:            "suiprivkey1qt2yx236gj92a304eh49xzfmmypvaqeu9vvcjnaz09eyaumg0466vw83env",
				toSuiAddress:            "0xd4b4dcbc801c2d465b3f7a44cb1acad9702c14f790302625fef20b0811dc636a",
				publicKeyToRawBytes:     []byte{2, 124, 18, 113, 171, 75, 241, 31, 57, 238, 45, 0, 213, 64, 91, 78, 8, 205, 136, 84, 150, 206, 253, 112, 171, 155, 195, 138, 172, 188, 252, 40, 20},
				publicKeyToSuiBytes:     []byte{2, 2, 124, 18, 113, 171, 75, 241, 31, 57, 238, 45, 0, 213, 64, 91, 78, 8, 205, 136, 84, 150, 206, 253, 112, 171, 155, 195, 138, 172, 188, 252, 40, 20},
				publicKeyToBase64:       "AnwScatL8R857i0A1UBbTgjNiFSWzv1wq5vDiqy8/CgU",
				publicKeyToSuiPublicKey: "AgJ8EnGrS/EfOe4tANVAW04IzYhUls79cKubw4qsvPwoFA==",
				publicKeyToSuiAddress:   "0xd4b4dcbc801c2d465b3f7a44cb1acad9702c14f790302625fef20b0811dc636a",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AkJXhYmR3t2UnTgeDJYyOBPaScS0+Aoip9ehG4DEy6BoBPMJIGUUyI4SHoPi4KSjaMCEDWIk1WUI/Ot7Q6nBCdgCfBJxq0vxHznuLQDVQFtOCM2IVJbO/XCrm8OKrLz8KBQ=",
			},
		},
		{
			path: "m/74'/784'/0'/0/100",
			expected: struct {
				getKeyScheme            string
				getSecretKey            string
				toSuiAddress            string
				publicKeyToRawBytes     []byte
				publicKeyToSuiBytes     []byte
				publicKeyToBase64       string
				publicKeyToSuiPublicKey string
				publicKeyToSuiAddress   string
				message                 string
				signature               string
			}{
				getKeyScheme:            "Secp256r1",
				getSecretKey:            "suiprivkey1qgw0mmwmjrnv0p5k4pz280pycp6kqf2el00xhtg7hm8dp9d8hjs2cw27tma",
				toSuiAddress:            "0xa4872fd99dc16d9374adee9759653be974b770bf437a47e7db7f3fc3edabfe37",
				publicKeyToRawBytes:     []byte{2, 172, 241, 252, 174, 88, 35, 149, 136, 100, 96, 254, 24, 99, 230, 135, 48, 246, 182, 211, 11, 245, 6, 98, 2, 29, 9, 75, 69, 59, 164, 117, 186},
				publicKeyToSuiBytes:     []byte{2, 2, 172, 241, 252, 174, 88, 35, 149, 136, 100, 96, 254, 24, 99, 230, 135, 48, 246, 182, 211, 11, 245, 6, 98, 2, 29, 9, 75, 69, 59, 164, 117, 186},
				publicKeyToBase64:       "Aqzx/K5YI5WIZGD+GGPmhzD2ttML9QZiAh0JS0U7pHW6",
				publicKeyToSuiPublicKey: "AgKs8fyuWCOViGRg/hhj5ocw9rbTC/UGYgIdCUtFO6R1ug==",
				publicKeyToSuiAddress:   "0xa4872fd99dc16d9374adee9759653be974b770bf437a47e7db7f3fc3edabfe37",

				message:   "EUhlbGxvIEdvIE1vZHVsZXMh",
				signature: "AnTTdvH53Oe3bwoocWFhFymLBFCch2BuGVIbPFr7THvGVcwH2mvFL9fx+kLHb43gAkbJVhaR/rqa6pVgz4Aa1ZQCrPH8rlgjlYhkYP4YY+aHMPa20wv1BmICHQlLRTukdbo=",
			},
		},
	}
	for _, test := range testDatas {
		t.Run(test.path, func(t *testing.T) {
			// For keypair
			keypair, err := secp256r1.DeriveKeypair(mnemonic, test.path)
			if err != nil {
				t.Fatalf("failed to derive secp256r1 keypair, msg: %v", err)
			}

			getKeyScheme := keypair.GetKeyScheme()
			if !reflect.DeepEqual(getKeyScheme, test.expected.getKeyScheme) {
				t.Errorf("unable to get key scheme, expected %s, but got %s", test.expected.getKeyScheme, getKeyScheme)
			}

			getSecretKey, err := keypair.GetSecretKey()
			if err != nil {
				t.Fatalf("failed to get secret key, msg: %v", err)
			}
			if !reflect.DeepEqual(getSecretKey, test.expected.getSecretKey) {
				t.Errorf("unable to get secret key, expected %v, but got %v", test.expected.getSecretKey, getSecretKey)
			}

			if !reflect.DeepEqual(keypair.ToSuiAddress(), test.expected.toSuiAddress) {
				t.Errorf("unable to sui address, expected %v, but got %v", test.expected.toSuiAddress, keypair.ToSuiAddress())
			}

			// For public key
			pubkey, err := keypair.GetPublicKey()
			if err != nil {
				t.Fatalf("failed to get public key, msg: %v", err)
			}
			if !bytes.Equal(pubkey.ToRawBytes(), test.expected.publicKeyToRawBytes) {
				t.Errorf("unable to get public key to raw bytes, expected %v, but got %v", test.expected.publicKeyToRawBytes, pubkey.ToRawBytes())
			}

			if !bytes.Equal(pubkey.ToSuiBytes(), test.expected.publicKeyToSuiBytes) {
				t.Errorf("unable to get public key to sui bytes, expected %v, but got %v", test.expected.publicKeyToSuiBytes, pubkey.ToSuiBytes())
			}

			if !reflect.DeepEqual(pubkey.ToBase64(), test.expected.publicKeyToBase64) {
				t.Errorf("unable to get public key to base64, expected %v, but got %v", test.expected.publicKeyToBase64, pubkey.ToBase64())
			}

			if !reflect.DeepEqual(pubkey.ToSuiPublicKey(), test.expected.publicKeyToSuiPublicKey) {
				t.Errorf("unable to get public key to sui public key, expected %v, but got %v", test.expected.publicKeyToSuiPublicKey, pubkey.ToSuiPublicKey())
			}

			if !reflect.DeepEqual(pubkey.ToSuiAddress(), test.expected.publicKeyToSuiAddress) {
				t.Errorf("unable to get public key to sui address, expected %v, but got %v", test.expected.publicKeyToSuiAddress, pubkey.ToSuiAddress())
			}

			// For signature
			data, err := keypair.SignPersonalMessage([]byte(message))
			if err != nil {
				t.Fatalf("keypair failed to sign personal message, msg: %v", err)
			}

			if !reflect.DeepEqual(data.Bytes, test.expected.message) {
				t.Errorf("unable to sign personal message, bytes expected %s, but got %v", test.expected.message, data.Bytes)
			}

			if !reflect.DeepEqual(data.Signature, test.expected.signature) {
				t.Errorf("unable to sign personal message, signature expected %s, but got %v", test.expected.signature, data.Signature)
			}

			pass, err := pubkey.VerifyPersonalMessage([]byte(message), data.Signature)
			if err != nil {
				t.Fatalf("public key failed to verify personal message, msg: %v", err)
			}
			if !pass {
				t.Errorf("unable to verify personal message, expected %v, but got %v", true, pass)
			}
		})
	}
}

func TestMultisig(t *testing.T) {
	testDatas := []struct {
		ed25519Path     string
		ed25519Weight   uint8
		secp256k1Path   string
		secp256k1Weight uint8
		secp256r1Path   string
		secp256r1Weight uint8
		message         string
		threshold       uint16
		expected        struct {
			flag               uint8
			threshold          uint16
			toBase64           string
			toRawBytes         []byte
			toSuiPublicKey     string
			toSuiBytes         []byte
			toSuiAddress       string
			ed25519Signature   string
			secp256k1Signature string
			secp256r1Signature string
			combineSignature   string
		}
	}{
		{
			ed25519Path:     "m/44'/784'/0'/0'/100'",
			ed25519Weight:   2,
			secp256k1Path:   "m/54'/784'/0'/0/100",
			secp256k1Weight: 2,
			secp256r1Path:   "m/74'/784'/0'/0/100",
			secp256r1Weight: 2,
			message:         "Hello Sui MultiSig!",
			threshold:       3,
			expected: struct {
				flag               uint8
				threshold          uint16
				toBase64           string
				toRawBytes         []byte
				toSuiPublicKey     string
				toSuiBytes         []byte
				toSuiAddress       string
				ed25519Signature   string
				secp256k1Signature string
				secp256r1Signature string
				combineSignature   string
			}{
				flag:               3,
				threshold:          3,
				toBase64:           "AwC1IeWE/ON0Ya62CCt6T3e7pOVmAqPosENNfuwr/hifXAIBAinQ0sYDdCCxz6RU5cV6BsVTayq5GLCRFZU8IK+sxHIhAgICrPH8rlgjlYhkYP4YY+aHMPa20wv1BmICHQlLRTukdboCAwA=",
				toRawBytes:         []byte{3, 0, 181, 33, 229, 132, 252, 227, 116, 97, 174, 182, 8, 43, 122, 79, 119, 187, 164, 229, 102, 2, 163, 232, 176, 67, 77, 126, 236, 43, 254, 24, 159, 92, 2, 1, 2, 41, 208, 210, 198, 3, 116, 32, 177, 207, 164, 84, 229, 197, 122, 6, 197, 83, 107, 42, 185, 24, 176, 145, 21, 149, 60, 32, 175, 172, 196, 114, 33, 2, 2, 2, 172, 241, 252, 174, 88, 35, 149, 136, 100, 96, 254, 24, 99, 230, 135, 48, 246, 182, 211, 11, 245, 6, 98, 2, 29, 9, 75, 69, 59, 164, 117, 186, 2, 3, 0},
				toSuiPublicKey:     "AwMAtSHlhPzjdGGutggrek93u6TlZgKj6LBDTX7sK/4Yn1wCAQIp0NLGA3Qgsc+kVOXFegbFU2squRiwkRWVPCCvrMRyIQICAqzx/K5YI5WIZGD+GGPmhzD2ttML9QZiAh0JS0U7pHW6AgMA",
				toSuiBytes:         []byte{3, 3, 0, 181, 33, 229, 132, 252, 227, 116, 97, 174, 182, 8, 43, 122, 79, 119, 187, 164, 229, 102, 2, 163, 232, 176, 67, 77, 126, 236, 43, 254, 24, 159, 92, 2, 1, 2, 41, 208, 210, 198, 3, 116, 32, 177, 207, 164, 84, 229, 197, 122, 6, 197, 83, 107, 42, 185, 24, 176, 145, 21, 149, 60, 32, 175, 172, 196, 114, 33, 2, 2, 2, 172, 241, 252, 174, 88, 35, 149, 136, 100, 96, 254, 24, 99, 230, 135, 48, 246, 182, 211, 11, 245, 6, 98, 2, 29, 9, 75, 69, 59, 164, 117, 186, 2, 3, 0},
				toSuiAddress:       "0xfb061fbd2807b3b6f74cc8f238141ee8e72623a1cd35b9cd6f6bce40be9caccb",
				ed25519Signature:   "AIgT8hGxxolm+7DH2GYa0F8GH29su6hUXKESLTLO58O7FbE32o98EVHmfMGhBEDL3i2crqZg6EZoLM1yAPst/Qy1IeWE/ON0Ya62CCt6T3e7pOVmAqPosENNfuwr/hifXA==",
				secp256k1Signature: "AUCEfb7JJL6T3U2ELhxYGxRE/ys1lnVg9SBOT5dYNIapaE58rOnch+nLlPoSaAqNL0Gnn+wvYvGsO5mGG/u++pgCKdDSxgN0ILHPpFTlxXoGxVNrKrkYsJEVlTwgr6zEciE=",
				secp256r1Signature: "Ak3cCty+GNPSMkwKEb4UwxJv8KAxO/4sXGWGIT6K/NsnASzh/LgJFDturYbfs7eJUmY1QrdQCsx4JS7sL6UCF64CrPH8rlgjlYhkYP4YY+aHMPa20wv1BmICHQlLRTukdbo=",
				combineSignature:   "AwMAiBPyEbHGiWb7sMfYZhrQXwYfb2y7qFRcoRItMs7nw7sVsTfaj3wRUeZ8waEEQMveLZyupmDoRmgszXIA+y39DAFAhH2+ySS+k91NhC4cWBsURP8rNZZ1YPUgTk+XWDSGqWhOfKzp3Ifpy5T6EmgKjS9Bp5/sL2LxrDuZhhv7vvqYAk3cCty+GNPSMkwKEb4UwxJv8KAxO/4sXGWGIT6K/NsnASzh/LgJFDturYbfs7eJUmY1QrdQCsx4JS7sL6UCF64HAAMAtSHlhPzjdGGutggrek93u6TlZgKj6LBDTX7sK/4Yn1wCAQIp0NLGA3Qgsc+kVOXFegbFU2squRiwkRWVPCCvrMRyIQICAqzx/K5YI5WIZGD+GGPmhzD2ttML9QZiAh0JS0U7pHW6AgMA",
			},
		},
		{
			ed25519Path:     "m/44'/784'/0'/0'/0'",
			ed25519Weight:   1,
			secp256k1Path:   "m/54'/784'/0'/0/0",
			secp256k1Weight: 1,
			secp256r1Path:   "m/74'/784'/0'/0/0",
			secp256r1Weight: 1,
			message:         "Hello Sui MultiSig!",
			threshold:       2,
			expected: struct {
				flag               uint8
				threshold          uint16
				toBase64           string
				toRawBytes         []byte
				toSuiPublicKey     string
				toSuiBytes         []byte
				toSuiAddress       string
				ed25519Signature   string
				secp256k1Signature string
				secp256r1Signature string
				combineSignature   string
			}{
				flag:               3,
				threshold:          2,
				toBase64:           "AwCQC02B7s6j3y90sUIAxPTPP0mvrKemNP/Sz2/4K9rs8gEBAmI9hg9GzOkRfT8aw4K3nFmSigBKGYZWGpnfKoUWfPWFAQIDQBm8qKh4RYpj5b9T8whV4xBw97V8+dzyZcmL2xe7F8QBAgA=",
				toRawBytes:         []byte{3, 0, 144, 11, 77, 129, 238, 206, 163, 223, 47, 116, 177, 66, 0, 196, 244, 207, 63, 73, 175, 172, 167, 166, 52, 255, 210, 207, 111, 248, 43, 218, 236, 242, 1, 1, 2, 98, 61, 134, 15, 70, 204, 233, 17, 125, 63, 26, 195, 130, 183, 156, 89, 146, 138, 0, 74, 25, 134, 86, 26, 153, 223, 42, 133, 22, 124, 245, 133, 1, 2, 3, 64, 25, 188, 168, 168, 120, 69, 138, 99, 229, 191, 83, 243, 8, 85, 227, 16, 112, 247, 181, 124, 249, 220, 242, 101, 201, 139, 219, 23, 187, 23, 196, 1, 2, 0},
				toSuiPublicKey:     "AwMAkAtNge7Oo98vdLFCAMT0zz9Jr6ynpjT/0s9v+Cva7PIBAQJiPYYPRszpEX0/GsOCt5xZkooAShmGVhqZ3yqFFnz1hQECA0AZvKioeEWKY+W/U/MIVeMQcPe1fPnc8mXJi9sXuxfEAQIA",
				toSuiBytes:         []byte{3, 3, 0, 144, 11, 77, 129, 238, 206, 163, 223, 47, 116, 177, 66, 0, 196, 244, 207, 63, 73, 175, 172, 167, 166, 52, 255, 210, 207, 111, 248, 43, 218, 236, 242, 1, 1, 2, 98, 61, 134, 15, 70, 204, 233, 17, 125, 63, 26, 195, 130, 183, 156, 89, 146, 138, 0, 74, 25, 134, 86, 26, 153, 223, 42, 133, 22, 124, 245, 133, 1, 2, 3, 64, 25, 188, 168, 168, 120, 69, 138, 99, 229, 191, 83, 243, 8, 85, 227, 16, 112, 247, 181, 124, 249, 220, 242, 101, 201, 139, 219, 23, 187, 23, 196, 1, 2, 0},
				toSuiAddress:       "0x357ef6da9a57094cb90b1d6c40a5f0a5c533e48ed69d02f66e292ceff2f80cd9",
				ed25519Signature:   "AIzXgqJN8Rwqt22Ap2pwkx9+ucchpDeAsCiAOZaOZZqbgqaEybaIXjys7L4mojDMo+3pWlSeLpV6WJN1yO7tVQyQC02B7s6j3y90sUIAxPTPP0mvrKemNP/Sz2/4K9rs8g==",
				secp256k1Signature: "AdDfYvAb/NGpRNcvnItMof1rY2rE8c8QPcpdebS3FmrRNGePO+AasY1Eklhe50P0KzVvLHK4PGgzyb2PpJDOQYoCYj2GD0bM6RF9PxrDgrecWZKKAEoZhlYamd8qhRZ89YU=",
				secp256r1Signature: "AohHLG9bUkjYFyXxIiYY7JprgAnerwiPYX1HA8Byqug3aaTMcqeWNYJu/s5Z2e0256oZprIttsZ36YBpc381sa4DQBm8qKh4RYpj5b9T8whV4xBw97V8+dzyZcmL2xe7F8Q=",
				combineSignature:   "AwMAjNeCok3xHCq3bYCnanCTH365xyGkN4CwKIA5lo5lmpuCpoTJtohePKzsviaiMMyj7elaVJ4ulXpYk3XI7u1VDAHQ32LwG/zRqUTXL5yLTKH9a2NqxPHPED3KXXm0txZq0TRnjzvgGrGNRJJYXudD9Cs1byxyuDxoM8m9j6SQzkGKAohHLG9bUkjYFyXxIiYY7JprgAnerwiPYX1HA8Byqug3aaTMcqeWNYJu/s5Z2e0256oZprIttsZ36YBpc381sa4HAAMAkAtNge7Oo98vdLFCAMT0zz9Jr6ynpjT/0s9v+Cva7PIBAQJiPYYPRszpEX0/GsOCt5xZkooAShmGVhqZ3yqFFnz1hQECA0AZvKioeEWKY+W/U/MIVeMQcPe1fPnc8mXJi9sXuxfEAQIA",
			},
		},
	}

	for _, test := range testDatas {
		t.Run(test.expected.toSuiAddress, func(t *testing.T) {
			ed25519Keypair, err := ed25519.DeriveKeypair(mnemonic, test.ed25519Path)
			if err != nil {
				t.Fatalf("failed to derive ed25519 keypair, msg: %v", err)
			}

			secp256k1Keypair, err := secp256k1.DeriveKeypair(mnemonic, test.secp256k1Path)
			if err != nil {
				t.Fatalf("failed to derive secp256k1 keypair, msg: %v", err)
			}

			secp256r1Keypair, err := secp256r1.DeriveKeypair(mnemonic, test.secp256r1Path)
			if err != nil {
				t.Fatalf("failed to derive secp256r1 keypair, msg: %v", err)
			}

			ed25519PublicKey, err := ed25519Keypair.GetPublicKey()
			if err != nil {
				t.Fatalf("failed to get ed25519 public key, msg: %v", err)
			}

			secp256k1PublicKey, err := secp256k1Keypair.GetPublicKey()
			if err != nil {
				t.Fatalf("failed to get secp256k1 public key, msg: %v", err)
			}

			secp256r1PublicKey, err := secp256r1Keypair.GetPublicKey()
			if err != nil {
				t.Fatalf("failed to get secp256r1 public key, msg: %v", err)
			}

			multisigPubKey, err := new(multisig.MultiSigPublicKey).FromPublicKeys([]multisig.PublicKeyWeightPair{
				{PublicKey: ed25519PublicKey, Weight: test.ed25519Weight},
				{PublicKey: secp256k1PublicKey, Weight: test.secp256k1Weight},
				{PublicKey: secp256r1PublicKey, Weight: test.secp256r1Weight},
			}, test.threshold)
			if err != nil {
				t.Fatalf("failed to new multisig public key, msg: %v", err)
			}

			if !reflect.DeepEqual(test.expected.flag, multisigPubKey.Flag()) {
				t.Errorf("expected flag %v, but got %v", test.expected.flag, multisigPubKey.Flag())
			}
			if !reflect.DeepEqual(test.expected.threshold, multisigPubKey.GetThreshold()) {
				t.Errorf("expected threshold %v, but got %v", test.expected.threshold, multisigPubKey.GetThreshold())
			}
			if !reflect.DeepEqual(test.expected.toBase64, multisigPubKey.ToBase64()) {
				t.Errorf("expected base64 %v, but got %v", test.expected.toBase64, multisigPubKey.ToBase64())
			}
			if !reflect.DeepEqual(test.expected.toSuiPublicKey, multisigPubKey.ToSuiPublicKey()) {
				t.Errorf("expected sui public key %v, but got %v", test.expected.toSuiPublicKey, multisigPubKey.ToSuiPublicKey())
			}
			if !reflect.DeepEqual(test.expected.toSuiAddress, multisigPubKey.ToSuiAddress()) {
				t.Errorf("expected sui address %v, but got %v", test.expected.toSuiAddress, multisigPubKey.ToSuiAddress())
			}
			if !bytes.Equal(test.expected.toRawBytes, multisigPubKey.ToRawBytes()) {
				t.Errorf("expected raw bytes %v, but got %v", test.expected.toRawBytes, multisigPubKey.ToRawBytes())
			}
			if !bytes.Equal(test.expected.toSuiBytes, multisigPubKey.ToSuiBytes()) {
				t.Errorf("expected raw sui bytes %v, but got %v", test.expected.toSuiBytes, multisigPubKey.ToSuiBytes())
			}

			m, err := multisig.NewMultiSigPublicKey(multisigPubKey.ToRawBytes())
			if err != nil {
				t.Fatalf("failed to new multisig public key, msg: %v", err)
			}

			if !bytes.Equal(m.ToRawBytes(), multisigPubKey.ToRawBytes()) {
				t.Errorf("expected raw bytes %v, but got %v", multisigPubKey.ToRawBytes(), m.ToRawBytes())
			}

			signatureData1, err := ed25519Keypair.SignPersonalMessage([]byte(test.message))
			if err != nil {
				t.Fatalf("failed to sign ed25519 personal message, msg: %v", err)
			}
			if !reflect.DeepEqual(signatureData1.Signature, test.expected.ed25519Signature) {
				t.Errorf("expected ed25519 signature %v, but got %v", test.expected.ed25519Signature, signatureData1.Signature)
			}

			signatureData2, err := secp256k1Keypair.SignPersonalMessage([]byte(test.message))
			if err != nil {
				t.Fatalf("failed to sign secp256k1 personal message, msg: %v", err)
			}
			if !reflect.DeepEqual(signatureData2.Signature, test.expected.secp256k1Signature) {
				t.Errorf("expected secp256k1 signature %v, but got %v", test.expected.secp256k1Signature, signatureData2.Signature)
			}

			signatureData3, err := secp256r1Keypair.SignPersonalMessage([]byte(test.message))
			if err != nil {
				t.Fatalf("failed to sign secp256r1 personal message, msg: %v", err)
			}
			if !reflect.DeepEqual(signatureData3.Signature, test.expected.secp256r1Signature) {
				t.Errorf("expected secp256r1 signature %v, but got %v", test.expected.secp256r1Signature, signatureData3.Signature)
			}

			signature, err := multisigPubKey.CombinePartialSignatures([]cryptography.SerializedSignature{signatureData1.Signature, signatureData2.Signature, signatureData3.Signature})
			if err != nil {
				t.Fatalf("failed to combine partial signatures, msg: %v", err)
			}
			if !reflect.DeepEqual(signature, test.expected.combineSignature) {
				t.Errorf("expected signature %v, but got %v", test.expected.combineSignature, signature)
			}
		})

	}
}
