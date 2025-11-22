package encryption

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		plainText := []byte("hello world")
		cipherText, err := Encrypt(plainText, key)
		require.NoError(t, err)
		assert.NotEmpty(t, cipherText)
		assert.NotEqual(t, plainText, cipherText)

		decrypted, err := Decrypt(cipherText, key)
		require.NoError(t, err)
		assert.Equal(t, plainText, decrypted)
	})

	t.Run("invalid key length", func(t *testing.T) {
		_, err := Encrypt([]byte("test"), []byte("short"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encryption key must be 32 bytes")

		_, err = Decrypt([]byte("test"), []byte("short"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decryption key must be 32 bytes")
	})

	t.Run("invalid ciphertext", func(t *testing.T) {
		_, err := Decrypt([]byte("short"), key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid encrypted data")
	})

	t.Run("decrypt with wrong key", func(t *testing.T) {
		plainText := []byte("hello world")
		cipherText, err := Encrypt(plainText, key)
		require.NoError(t, err)

		wrongKey := make([]byte, 32)
		_, err = rand.Read(wrongKey)
		require.NoError(t, err)

		_, err = Decrypt(cipherText, wrongKey)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decrypt data (potentially incorrect key)")
	})
}

func TestEncryptDecryptString(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	t.Run("base64 success", func(t *testing.T) {
		plainText := "hello world"
		cipherText, err := EncryptString(plainText, key, EncodingBase64)
		require.NoError(t, err)
		assert.NotEmpty(t, cipherText)

		decrypted, err := DecryptString(cipherText, key, EncodingBase64)
		require.NoError(t, err)
		assert.Equal(t, plainText, decrypted)
	})

	t.Run("hex success", func(t *testing.T) {
		plainText := "hello world"
		cipherText, err := EncryptString(plainText, key, EncodingHex)
		require.NoError(t, err)
		assert.NotEmpty(t, cipherText)

		decrypted, err := DecryptString(cipherText, key, EncodingHex)
		require.NoError(t, err)
		assert.Equal(t, plainText, decrypted)
	})

	t.Run("invalid encoding", func(t *testing.T) {
		_, err := EncryptString("test", key, "invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid cipher text encoding")

		_, err = DecryptString("test", key, "invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid cipher text encoding")
	})

	t.Run("invalid base64 ciphertext", func(t *testing.T) {
		_, err := DecryptString("invalid-base64", key, EncodingBase64)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error during decode base64")
	})

	t.Run("invalid hex ciphertext", func(t *testing.T) {
		_, err := DecryptString("invalid-hex", key, EncodingHex)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error during decode hex")
	})
}
