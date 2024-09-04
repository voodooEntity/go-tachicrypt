package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// deriveKey converts a password string into a 32-byte key using SHA-256
func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func EncryptWithRandomKey(plaintext string) (string, string, error) {
	// Generate a random key
	key := make([]byte, aes.BlockSize)
	if _, err := rand.Read(key); err != nil {
		return "", "", errors.New("error generating random key")
	}

	// Encrypt using the generated key
	ciphertext, err := EncryptWithPassword(plaintext, string(key))
	if err != nil {
		return "", "", err
	}

	// Encode the key in Base64 for storage or transmission
	encodedKey := base64.StdEncoding.EncodeToString(key)

	return ciphertext, encodedKey, nil
}

// DecryptWithRandomKey decrypts a ciphertext using AES with a randomly generated key.
func DecryptWithRandomKey(ciphertext string, encodedKey string) (string, error) {
	// Decode the key
	keyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return "", errors.New("error decoding key")
	}

	// Decrypt using the decoded key
	return DecryptWithPassword(ciphertext, string(keyBytes))
}

func DecryptWithPassword(ciphertext string, password string) (string, error) {
	// Decode the ciphertext and IV
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New("error decoding ciphertext")
	}

	// Split the ciphertext into IV and actual ciphertext
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	// Derive a 32-byte key from the password
	key := deriveKey(password)
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %w", err)
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", fmt.Errorf("error creating GCM cipher: %w", err)
	}

	// Decrypt the ciphertext using GCM
	plaintext, err := gcm.Open(nil, iv, ciphertextBytes, nil)
	if err != nil {
		return "", errors.New("error decrypting ciphertext")
	}

	return string(plaintext), nil
}

func EncryptWithPassword(plaintext string, password string) (string, error) {
	// Derive a 32-byte key from the password
	key := deriveKey(password)
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %+w", err)
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", fmt.Errorf("error creating GCM cipher: %+w", err)
	}

	// Generate a random initialization vector (IV)
	iv := make([]byte, gcm.NonceSize()) // gcm.NonceSize() typically returns 12
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Encrypt the plaintext using GCM
	ciphertext := gcm.Seal(nil, iv, []byte(plaintext), nil)

	// Encode the ciphertext and IV in Base64 for storage or transmission
	return base64.StdEncoding.EncodeToString(append(iv, ciphertext...)), nil
}
