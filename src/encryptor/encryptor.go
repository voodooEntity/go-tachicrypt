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

// test hooks / indirection for easier unit testing of error paths
var (
    randReader     = rand.Reader
    aesNewCipher   = aes.NewCipher
    cipherNewGCM   = cipher.NewGCM
)

// deriveKey converts a password string into a 32-byte key using SHA-256
func deriveKey(password string) []byte {
    hash := sha256.Sum256([]byte(password))
    return hash[:]
}

func EncryptWithRandomKey(data []byte) ([]byte, string, error) {
	// Generate a random key
	key := make([]byte, aes.BlockSize)
 if _, err := randReader.Read(key); err != nil {
        return []byte{}, "", errors.New("error generating random key")
    }

	// Encrypt using the generated key
	ciphertextBytes, err := EncryptWithPassword(data, string(key))
	if err != nil {
		return []byte{}, "", err
	}

	// Encode the key in Base64 for storage or transmission
	encodedKey := base64.StdEncoding.EncodeToString(key)

	return ciphertextBytes, encodedKey, nil
}

// DecryptWithRandomKey decrypts a ciphertext using AES with a randomly generated key.
func DecryptWithRandomKey(ciphertextBytes []byte, encodedKey string) ([]byte, error) {
	// Decode the key
	keyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return []byte{}, fmt.Errorf("error decoding key: %+w ", err)
	}

	// Decrypt using the decoded key
	return DecryptWithPassword(ciphertextBytes, string(keyBytes))
}

func DecryptWithPassword(ciphertextBytes []byte, password string) ([]byte, error) {
	// Derive a 32-byte key from the password
	key := deriveKey(password)
 aesCipher, err := aesNewCipher(key)
 if err != nil {
     return []byte{}, fmt.Errorf("error creating AES cipher: %w", err)
 }

 // Create a new GCM cipher
 gcm, err := cipherNewGCM(aesCipher)
 if err != nil {
     return []byte{}, fmt.Errorf("error creating GCM cipher: %w", err)
 }

	// Split the ciphertext into IV and actual ciphertext
	iv := ciphertextBytes[:gcm.NonceSize()]
	ciphertextBytes = ciphertextBytes[gcm.NonceSize():]

	// Decrypt the ciphertext using GCM
	dst := []byte{}
	plaintext, err := gcm.Open(dst, iv, ciphertextBytes, nil)
	if err != nil {
		return []byte{}, errors.New("error decrypting ciphertext")
	}

	return plaintext, nil
}

func EncryptWithPassword(data []byte, password string) ([]byte, error) {
	// Derive a 32-byte key from the password
	key := deriveKey(password)
 aesCipher, err := aesNewCipher(key)
 if err != nil {
     return []byte{}, fmt.Errorf("error creating AES cipher: %+w", err)
 }

 // Create a new GCM cipher
 gcm, err := cipherNewGCM(aesCipher)
 if err != nil {
     return []byte{}, fmt.Errorf("error creating GCM cipher: %+w", err)
 }

	// Generate a random initialization vector (IV)
	iv := make([]byte, gcm.NonceSize()) // gcm.NonceSize() typically returns 12
 if _, err := io.ReadFull(randReader, iv); err != nil {
        panic(err)
    }

	// Encrypt the plaintext using GCM
	ciphertext := gcm.Seal(nil, iv, data, nil)

	// Encode the ciphertext and IV in Base64 for storage or transmission
	return append(iv, ciphertext...), nil
}
