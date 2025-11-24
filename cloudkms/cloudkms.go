package cloudkms

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/cockroachdb/errors"
)

type Encoding string

const (
	EncodingBase64 Encoding = "base64"
	EncodingHex    Encoding = "hex"
)

type Client interface {
	// Encrypt encrypts the given plaintext bytes
	Encrypt(ctx context.Context, plainText []byte) ([]byte, error)

	// Decrypt decrypts the given ciphertext bytes
	Decrypt(ctx context.Context, cipherText []byte) ([]byte, error)

	// EncryptString encrypts the given plaintext and returns the ciphertext as a string.
	EncryptString(ctx context.Context, plainText string, cipherTextEncoding Encoding) (string, error)

	// DecryptString decrypts the given ciphertext and returns the plaintext as a string.
	DecryptString(ctx context.Context, cipherText string, cipherTextEncoding Encoding) (string, error)

	Close() error
}

var _ Client = (*client)(nil)

type client struct {
	kmsClient *kms.KeyManagementClient
	keyName   string
}

type Config struct {
	Project  string `env:"PROJECT" mapstructure:"project"`
	Location string `env:"LOCATION" mapstructure:"location"`
	KeyRing  string `env:"KEYRING" mapstructure:"keyring"`
	Key      string `env:"KEY" mapstructure:"key"`
}

func New(ctx context.Context, config Config) (Client, error) {
	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	keyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", config.Project, config.Location, config.KeyRing, config.Key)
	return &client{kmsClient: kmsClient, keyName: keyName}, nil
}

func (c *client) Encrypt(ctx context.Context, plainText []byte) ([]byte, error) {
	encryptRequest := &kmspb.EncryptRequest{
		Name:      c.keyName,
		Plaintext: plainText,
	}
	encryptResponse, err := c.kmsClient.Encrypt(ctx, encryptRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error during encrypt plaintext")
	}
	return encryptResponse.Ciphertext, nil
}

func (c *client) Decrypt(ctx context.Context, cipherText []byte) ([]byte, error) {
	decryptRequest := &kmspb.DecryptRequest{
		Name:       c.keyName,
		Ciphertext: cipherText,
	}
	decryptResponse, err := c.kmsClient.Decrypt(ctx, decryptRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error during decrypt ciphertext")
	}
	return decryptResponse.Plaintext, nil
}

func (c *client) EncryptString(ctx context.Context, plainText string, cipherTextEncoding Encoding) (string, error) {
	encryptedBytes, err := c.Encrypt(ctx, []byte(plainText))
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

func (c *client) DecryptString(ctx context.Context, cipherText string, cipherTextEncoding Encoding) (string, error) {
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

	decryptedBytes, err := c.Decrypt(ctx, cipherTextBytes)
	if err != nil {
		return "", errors.Wrap(err, "error during decrypt string")
	}
	return string(decryptedBytes), nil
}

func (c *client) Close() error {
	err := c.kmsClient.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
