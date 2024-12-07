package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/random"
	"github.com/aiechoic/admin/core/viper"
	"io"
)

const DefaultAESConfig = "aes-gcm"

var initAESConfig = `
# AES-GCM 加密算法配置

# AES-GCM 加密算法的密钥, 长度必须是 16, 24, 32 字节
# 生产环境中请使用环境变量来设置密钥, 如: export AES-GCM_KEY=your_key
key: "` + random.String(32) + `"
`

type aesConfig struct {
	Key string `mapstructure:"key"`
}

var AESProviders = ioc.NewProviders[*AESCipher](func(name string, args ...any) *ioc.Provider[*AESCipher] {
	return ioc.NewProvider(func(c *ioc.Container) (*AESCipher, error) {
		vp, err := viper.GetViper(name, initAESConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg aesConfig
		if err := vp.Unmarshal(&cfg); err != nil {
			return nil, err
		}
		return NewAESCipher([]byte(cfg.Key))
	})
})

func GetAESCipher(name string, c *ioc.Container) (*AESCipher, error) {
	return AESProviders.GetProvider(name).Get(c)
}

func GetDefaultAESCipher(c *ioc.Container) (*AESCipher, error) {
	return GetAESCipher(DefaultAESConfig, c)
}

// AESCipher 是 AES-GCM 加密算法的封装
type AESCipher struct {
	aesGCM cipher.AEAD // AES-GCM 模式实例
}

// NewAESCipher 创建一个 AES-GCM 加密算法实例
func NewAESCipher(key []byte) (*AESCipher, error) {
	// 检查密钥长度是否合法（16, 24, 32 字节分别对应 AES-128, AES-192, AES-256）
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: key must be 16, 24, or 32 bytes")
	}

	// 初始化 AES 加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建 GCM 模式实例
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESCipher{
		aesGCM: aesGCM,
	}, nil
}

// Encrypt AES-GCM 加密
func (a *AESCipher) EncryptText(plainText []byte) ([]byte, error) {
	// 生成随机 nonce
	nonce := make([]byte, a.aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := a.aesGCM.Seal(nil, nonce, plainText, nil)

	// 拼接 nonce 和加密数据，并编码为 Base64
	encrypted := append(nonce, ciphertext...)
	return EncodeBase64Bytes(base64.StdEncoding, encrypted), nil
}

// Decrypt AES-GCM 解密
func (a *AESCipher) DecryptText(cipherText []byte) ([]byte, error) {
	// 解码 Base64 数据
	encrypted, err := DecodeBase64Bytes(base64.StdEncoding, cipherText)
	if err != nil {
		return nil, err
	}

	nonceSize := a.aesGCM.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, errors.New("encrypted data too short")
	}

	// 分离 nonce 和加密数据
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// 解密数据
	plaintext, err := a.aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (a *AESCipher) EncryptJSON(plainValue interface{}) (cipherText []byte, err error) {
	data, err := json.Marshal(plainValue)
	if err != nil {
		return nil, err
	}
	return a.EncryptText(data)
}

func (a *AESCipher) DecryptJSON(cipherText []byte, plainValue interface{}) error {
	data, err := a.DecryptText(cipherText)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, plainValue)
}
