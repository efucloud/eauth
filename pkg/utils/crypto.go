package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type PaddingMode string

const PKCS5 PaddingMode = "PKCS5"
const PKCS7 PaddingMode = "PKCS7"
const ZEROS PaddingMode = "ZEROS"

func Padding(padding PaddingMode, src []byte, blockSize int) []byte {
	switch padding {
	case PKCS5:
		src = PKCS5Padding(src, blockSize)
	case PKCS7:
		src = PKCS7Padding(src, blockSize)
	case ZEROS:
		src = ZerosPadding(src, blockSize)
	}
	return src
}

func UnPadding(padding PaddingMode, src []byte) ([]byte, error) {
	switch padding {
	case PKCS5:
		return PKCS5UnPadding(src)
	case PKCS7:
		return PKCS7UnPadding(src)
	case ZEROS:
		return ZerosUnPadding(src)
	}
	return src, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	return PKCS7Padding(src, blockSize)
}

func PKCS5UnPadding(src []byte) ([]byte, error) {
	return PKCS7UnPadding(src)
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return src, fmt.Errorf("src length is 0")
	}
	unpadding := int(src[length-1])
	if length < unpadding {
		return src, fmt.Errorf("src length is less than unpadding")
	}
	return src[:(length - unpadding)], nil
}

func ZerosPadding(src []byte, blockSize int) []byte {
	rem := len(src) % blockSize
	if rem == 0 {
		return src
	}
	return append(src, bytes.Repeat([]byte{0}, blockSize-rem)...)
}

func ZerosUnPadding(src []byte) ([]byte, error) {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1], nil
		}
	}
}

// AesSimpleEncrypt encrypts data with key using AES algorithm.
// In simple encryption mode, the user only needs to specify the key to complete the encryption.
// IV will be obtained by hashing the key. By default, PKCS7Padding and CBC modes are used.
// Return empty string if error occurs.
func AesSimpleEncrypt(data, key string) string {
	key = trimByMaxKeySize(key)
	keyBytes := ZerosPadding([]byte(key), aes.BlockSize)
	return AesCBCEncrypt(data, string(keyBytes), GenIVFromKey(key), PKCS7)
}

// AesCBCEncrypt encrypts data with key and iv using AES algorithm.
// You must make sure the length of key and iv is 16 bytes. This func tion does not perform any padding for key or iv.
func AesCBCEncrypt(data, key, iv string, paddingMode PaddingMode) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	src := Padding(paddingMode, []byte(data), block.BlockSize())
	encryptData := make([]byte, len(src))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(encryptData, src)
	return base64.StdEncoding.EncodeToString(encryptData)
}

// AesSimpleDecrypt decrypts data with key using AES algorithm.
// In simple decryption mode, the user only needs to specify the key to complete the decryption.
// This func tion will automatically obtain the IV by hashing the key.
func AesSimpleDecrypt(data, key string) string {
	key = trimByMaxKeySize(key)
	keyBytes := ZerosPadding([]byte(key), aes.BlockSize)
	return AesCBCDecrypt(data, string(keyBytes), GenIVFromKey(key), PKCS7)
}

// AesCBCDecrypt decrypts data with key and iv using AES algorithm.
func AesCBCDecrypt(data, key, iv string, paddingMode PaddingMode) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	decodeData, _ := base64.StdEncoding.DecodeString(data)
	decryptData := make([]byte, len(decodeData))
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(decryptData, decodeData)

	original, _ := UnPadding(paddingMode, decryptData)
	return string(original)
}

// GenIVFromKey generates IV from key.
func GenIVFromKey(key string) (iv string) {
	hashedKey := sha256.Sum256([]byte(key))
	return trimByBlockSize(hex.EncodeToString(hashedKey[:]))
}

func trimByBlockSize(key string) string {
	if len(key) > aes.BlockSize {
		return key[:aes.BlockSize]
	}
	return key
}

func trimByMaxKeySize(key string) string {
	if len(key) > 32 {
		return key[:32]
	}
	return key
}
