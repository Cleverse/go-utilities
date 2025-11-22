package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/cockroachdb/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

type Encoding string

const (
	EncodingBase64 Encoding = "base64"
	EncodingHex    Encoding = "hex"
)

func Encrypt(plainText []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher")
	}

	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "failed to generate nonce")
	}

	ciphertext := aead.Seal(nil, nonce, plainText, nil)

	encryptedData := append(nonce, ciphertext...)

	return encryptedData, nil
}

func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("decryption key must be 32 bytes")
	}

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher")
	}

	if len(cipherText) < chacha20poly1305.NonceSize {
		return nil, errors.New("invalid encrypted data")
	}
	nonce, ciphertext := cipherText[:chacha20poly1305.NonceSize], cipherText[chacha20poly1305.NonceSize:]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt data (potentially incorrect key)")
	}

	return plaintext, nil
}

func EncryptString(plainText string, key []byte, cipherTextEncoding Encoding) (string, error) {
	encryptedBytes, err := Encrypt([]byte(plainText), key)
	if err != nil {
		return "", errors.Wrap(err, "error during encrypt string")
	}
	switch cipherTextEncoding {
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(encryptedBytes), nil
	case EncodingHex:
		return hex.EncodeToString(encryptedBytes), nil
	default:
		return "", errors.Newf("invalid cipher text encoding: %s", cipherTextEncoding)
	}
}

func DecryptString(cipherText string, key []byte, cipherTextEncoding Encoding) (string, error) {
	var cipherTextBytes []byte
	var err error
	switch cipherTextEncoding {
	case EncodingBase64:
		cipherTextBytes, err = base64.StdEncoding.DecodeString(cipherText)
		if err != nil {
			return "", errors.Wrap(err, "error during decode base64")
		}
	case EncodingHex:
		cipherTextBytes, err = hex.DecodeString(cipherText)
		if err != nil {
			return "", errors.Wrap(err, "error during decode hex")
		}
	default:
		return "", errors.Newf("invalid cipher text encoding: %s", cipherTextEncoding)
	}

	decryptedBytes, err := Decrypt(cipherTextBytes, key)
	if err != nil {
		return "", errors.Wrap(err, "error during decrypt string")
	}
	return string(decryptedBytes), nil
}
