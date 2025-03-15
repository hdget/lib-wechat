package wxmp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/pkg/errors"
)

var (
	errAppIDNotMatch       = errors.New("app id not match")
	errInvalidBlockSize    = errors.New("invalid block size")
	errInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	errInvalidPKCS7Padding = errors.New("invalid padding on input")
)

// 解密加密信息获取微信用户信息
func decrypt(sessionKey, encryptedData, iv string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, errors.New("cipher.NewCBCDecrypter: IV length must equal block size")
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, errInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, errInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, errInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, errInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}
