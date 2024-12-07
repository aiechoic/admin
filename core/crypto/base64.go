package crypto

import "encoding/base64"

func DecodeBase64Bytes(enc *base64.Encoding, s []byte) ([]byte, error) {
	dst := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dst, s)
	return dst[:n], err
}

func EncodeBase64Bytes(enc *base64.Encoding, s []byte) []byte {
	dst := make([]byte, enc.EncodedLen(len(s)))
	enc.Encode(dst, s)
	return dst
}
