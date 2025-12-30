package encryptor

import (
    "crypto/cipher"
    "errors"
    "testing"
)

// stub types for aes/cipher factory errors
type badBlock struct{}

func (b *badBlock) BlockSize() int { return 16 }
func (b *badBlock) Encrypt(dst, src []byte) {}
func (b *badBlock) Decrypt(dst, src []byte) {}

func TestEncryptWithPassword_ErrorFromAesNewCipher(t *testing.T) {
    oldAES := aesNewCipher
    aesNewCipher = func(key []byte) (cipher.Block, error) {
        return nil, errors.New("aes fail")
    }
    defer func() { aesNewCipher = oldAES }()

    if _, err := EncryptWithPassword([]byte("x"), "pw"); err == nil {
        t.Fatalf("expected error from aesNewCipher")
    }
}

func TestEncryptWithPassword_ErrorFromCipherNewGCM(t *testing.T) {
    // Return a valid-looking block from aesNewCipher, then make cipherNewGCM fail
    oldAES := aesNewCipher
    aesNewCipher = func(key []byte) (cipher.Block, error) {
        return &badBlock{}, nil
    }
    oldGCM := cipherNewGCM
    cipherNewGCM = func(b cipher.Block) (cipher.AEAD, error) {
        return nil, errors.New("gcm fail")
    }
    defer func() { aesNewCipher = oldAES; cipherNewGCM = oldGCM }()

    if _, err := EncryptWithPassword([]byte("x"), "pw"); err == nil {
        t.Fatalf("expected error from cipherNewGCM in EncryptWithPassword")
    }
}

func TestDecryptWithPassword_ErrorFromAesNewCipher(t *testing.T) {
    oldAES := aesNewCipher
    aesNewCipher = func(key []byte) (cipher.Block, error) {
        return nil, errors.New("aes fail")
    }
    defer func() { aesNewCipher = oldAES }()

    if _, err := DecryptWithPassword([]byte{1,2,3,4,5,6,7,8,9,10,11,12}, "pw"); err == nil {
        t.Fatalf("expected error from aesNewCipher in decrypt")
    }
}

func TestDecryptWithPassword_ErrorFromCipherNewGCM(t *testing.T) {
    // Valid block
    oldAES := aesNewCipher
    aesNewCipher = func(key []byte) (cipher.Block, error) {
        return &badBlock{}, nil
    }
    oldGCM := cipherNewGCM
    cipherNewGCM = func(b cipher.Block) (cipher.AEAD, error) {
        return nil, errors.New("gcm fail")
    }
    defer func() { aesNewCipher = oldAES; cipherNewGCM = oldGCM }()

    if _, err := DecryptWithPassword([]byte{1,2,3,4,5,6,7,8,9,10,11,12}, "pw"); err == nil {
        t.Fatalf("expected error from cipherNewGCM in decrypt")
    }
}

func TestEncryptWithRandomKey_PropagatesEncryptError(t *testing.T) {
    // Force EncryptWithPassword to fail via cipherNewGCM hook
    oldAES := aesNewCipher
    aesNewCipher = func(key []byte) (cipher.Block, error) {
        return &badBlock{}, nil
    }
    oldGCM := cipherNewGCM
    cipherNewGCM = func(b cipher.Block) (cipher.AEAD, error) {
        return nil, errors.New("gcm fail")
    }
    // ensure key generation succeeds
    oldRand := randReader
    randReader = zeroReader{}
    defer func() { aesNewCipher = oldAES; cipherNewGCM = oldGCM; randReader = oldRand }()

    if _, _, err := EncryptWithRandomKey([]byte("data")); err == nil {
        t.Fatalf("expected error to propagate from EncryptWithPassword")
    }
}

// zeroReader returns zeros to avoid failing key generation in the above test
type zeroReader struct{}

func (z zeroReader) Read(p []byte) (int, error) { for i := range p { p[i] = 0 } ; return len(p), nil }
