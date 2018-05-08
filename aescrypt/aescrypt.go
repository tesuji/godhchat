// Package aescrypt implements an easy way to interact with aes cipher
package aescrypt

// Ref https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9
// or https://gist.github.com/manishtpatel/8222606
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AESCipher is an struct to interact with aes cipher
type AESCipher struct {
	bs  int      // blocksize
	key [32]byte // key to encrypt
}

// NewAESCipher inits a new AESCipher
func NewAESCipher(key [32]byte) *AESCipher {
	cip := new(AESCipher)
	cip.bs = 32
	cip.key = key
	return cip
}

// pad is the PKCS7 pad
func (cip *AESCipher) pad(src []byte) []byte {
	padding := cip.bs - len(src)%cip.bs
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// unpad is the PKCS7 unpad
func (cip *AESCipher) unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > cip.bs {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}
	padtext := bytes.Repeat([]byte{byte(unpadding)}, unpadding)
	if !bytes.Equal(src[length-unpadding:], padtext) {
		return src, nil
	}
	return src[:(length - unpadding)], nil
}

// Encrypt is an encryption using AES in CBC mode.
func (cip *AESCipher) Encrypt(text string) ([]byte, error) {
	block, err := aes.NewCipher(cip.key[:])
	if err != nil {
		return nil, err
	}

	msg := cip.pad([]byte(text))

	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], msg)

	dst := make([]byte, 4096)
	base64.URLEncoding.Encode(dst, ciphertext)
	return dst, nil
}

// Decrypt using AES in CBC mode. Expects the IV at the front of the string.
func (cip *AESCipher) Decrypt(text []byte) (string, error) {
	ciphertext := make([]byte, 4096)
	if _, err := base64.URLEncoding.Decode(ciphertext, text); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cip.key[:])
	if err != nil {
		return "", err
	}

	if (len(ciphertext) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	result, err := cip.unpad(ciphertext)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// func main() {
// 	key := []byte("LKHlhb899Y09olUi")
// 	encryptMsg, _ := encrypt(key, "Hello World")
// 	msg, _ := decrypt(key, encryptMsg)
// 	fmt.Println(msg) // Hello World
// }
