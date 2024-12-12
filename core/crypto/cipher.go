package crypto

type Cipher interface {
	EncryptText(plainText []byte) (cipherText []byte, err error)
	DecryptText(cipherText []byte) (plainText []byte, err error)
	EncryptJSON(plainValue interface{}) (cipherText []byte, err error)
	DecryptJSON(cipherText []byte, plainValue interface{}) error
}
