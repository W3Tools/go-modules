package ed25519_test

import (
	"reflect"
	"testing"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/cryptography"
	ed25519_keypair "github.com/W3Tools/go-modules/gmsui/keypairs/ed25519"
	"github.com/fardream/go-bcs/bcs"
)

func TestEd25519PublicKey(t *testing.T) {
	tests := []struct {
		suiAddress   string
		rawBytes     []byte
		publicKey    string
		suiBytes     []byte
		suiPublicKey string
		flag         uint8
		message      string
		signature    string
	}{
		{
			// m/44'/784'/0'/0'/0'
			suiAddress:   "0x5058cfdd208cf43ad09e95e064d16fec9033cddf424756b980bc2f645969ea5e",
			rawBytes:     []byte{254, 160, 35, 2, 213, 131, 174, 249, 17, 21, 83, 205, 86, 152, 194, 117, 175, 219, 42, 205, 48, 253, 187, 139, 108, 120, 141, 198, 49, 161, 180, 124},
			publicKey:    "/qAjAtWDrvkRFVPNVpjCda/bKs0w/buLbHiNxjGhtHw=",
			suiBytes:     []byte{0, 254, 160, 35, 2, 213, 131, 174, 249, 17, 21, 83, 205, 86, 152, 194, 117, 175, 219, 42, 205, 48, 253, 187, 139, 108, 120, 141, 198, 49, 161, 180, 124},
			suiPublicKey: "AP6gIwLVg675ERVTzVaYwnWv2yrNMP27i2x4jcYxobR8",
			flag:         0,
			message:      "hello",
			signature:    "AAX5Ny4lOGQuwCk6xomSNtuiFjmQNYf7FXxeKLJ7IX9vhQCWvDCO734lGCSJArZGGPXxanWzQunvHYkp6KbR6QL+oCMC1YOu+REVU81WmMJ1r9sqzTD9u4tseI3GMaG0fA==",
		},
		{
			// m/44'/784'/0'/0'/1'
			suiAddress:   "0xef65e3e07a2c4bc4370d7865ccdbbd9ae85962b6602525ef9cc68581203dc8a6",
			rawBytes:     []byte{117, 128, 142, 161, 104, 165, 131, 152, 47, 60, 74, 181, 173, 188, 9, 85, 36, 138, 213, 90, 57, 208, 249, 90, 80, 208, 100, 16, 96, 231, 208, 219},
			publicKey:    "dYCOoWilg5gvPEq1rbwJVSSK1Vo50PlaUNBkEGDn0Ns=",
			suiBytes:     []byte{0, 117, 128, 142, 161, 104, 165, 131, 152, 47, 60, 74, 181, 173, 188, 9, 85, 36, 138, 213, 90, 57, 208, 249, 90, 80, 208, 100, 16, 96, 231, 208, 219},
			suiPublicKey: "AHWAjqFopYOYLzxKta28CVUkitVaOdD5WlDQZBBg59Db",
			flag:         0,
			message:      "hello",
			signature:    "AOu/hsEv/4RqO9tezlmmn19net4x/cEqhpyETr3ZGA0I5uxagSQzPTw1oFdXH2EPcMZG44wX9gBeuZOKPPNmkQp1gI6haKWDmC88SrWtvAlVJIrVWjnQ+VpQ0GQQYOfQ2w==",
		},
		{
			// m/44'/784'/0'/0'/2'
			suiAddress:   "0xf7e1dfc2e9d82325dd5a6b67f683c6dceeb6f6fac7f78139614d5e80d9853b7f",
			rawBytes:     []byte{255, 141, 156, 70, 154, 124, 81, 87, 198, 68, 54, 241, 125, 144, 73, 233, 124, 253, 196, 64, 151, 221, 39, 157, 121, 56, 92, 238, 64, 163, 109, 57},
			publicKey:    "/42cRpp8UVfGRDbxfZBJ6Xz9xECX3SedeThc7kCjbTk=",
			suiBytes:     []byte{0, 255, 141, 156, 70, 154, 124, 81, 87, 198, 68, 54, 241, 125, 144, 73, 233, 124, 253, 196, 64, 151, 221, 39, 157, 121, 56, 92, 238, 64, 163, 109, 57},
			suiPublicKey: "AP+NnEaafFFXxkQ28X2QSel8/cRAl90nnXk4XO5Ao205",
			flag:         0,
			message:      "hello",
			signature:    "AJ0FNCly92/2jxEWnkp9hNKVR7vkdMqxCsXPZ3FlU5yvUPa57od6+1TmLH8LqTK1sNOe3caLh/q9VPv1Jz47FwL/jZxGmnxRV8ZENvF9kEnpfP3EQJfdJ515OFzuQKNtOQ==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.suiAddress, func(t *testing.T) {
			publicKey, err := ed25519_keypair.NewEd25519PublicKey(tt.publicKey)
			if err != nil {
				t.Fatalf("Unable to NewEd25519PublicKey, msg: %v", err)
			}

			if !reflect.DeepEqual(publicKey.ToSuiAddress(), tt.suiAddress) {
				t.Errorf("sui address expected %v, but got %v", tt.suiAddress, publicKey.ToSuiAddress())
			}

			if !gm.BytesEqual(publicKey.ToRawBytes(), tt.rawBytes) {
				t.Errorf("raw bytes expected %v, but got %v", tt.rawBytes, publicKey.ToRawBytes())
			}

			if !gm.BytesEqual(publicKey.ToSuiBytes(), tt.suiBytes) {
				t.Errorf("sui bytes expected %v, but got %v", tt.suiBytes, publicKey.ToSuiBytes())
			}

			if !reflect.DeepEqual(publicKey.ToBase64(), tt.publicKey) {
				t.Errorf("public key expected %v, but got %v", tt.suiPublicKey, publicKey.ToBase64())
			}

			if !reflect.DeepEqual(publicKey.ToSuiPublicKey(), tt.suiPublicKey) {
				t.Errorf("sui public key expected %v, but got %v", tt.suiPublicKey, publicKey.ToBase64())
			}

			if !reflect.DeepEqual(publicKey.Flag(), tt.flag) {
				t.Errorf("flag expected %v, but got %v", tt.flag, publicKey.Flag())
			}

			pass, err := publicKey.VerifyPersonalMessage([]byte(tt.message), tt.signature)
			if err != nil {
				t.Fatalf("unable to verify personal message, msg: %v", err)
			}
			if !pass {
				t.Errorf("verify personal message unpass")
			}

			bcsMessage, err := bcs.Marshal([]byte(tt.message))
			if err != nil {
				t.Fatalf("unable to marshal bcs, msg: %v", err)
			}
			pass, err = publicKey.VerifyWithIntent(bcsMessage, tt.signature, cryptography.PersonalMessage)
			if err != nil {
				t.Fatalf("unable to verify with intent, msg: %v", err)
			}
			if !pass {
				t.Errorf("verify with intent unpass")
			}
		})
	}
}
