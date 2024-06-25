package b64

import "encoding/base64"

func FromBase64(base64String string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64String)
}

func ToBase64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}
