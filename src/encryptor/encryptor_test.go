package encryptor

import (
    "bytes"
    "encoding/base64"
    "testing"
    "io"
)

func TestEncryptDecryptWithPassword_RoundTrip(t *testing.T) {
    data := []byte("some secret data for testing")
    pwd := "unit-test-password"

    ct, err := EncryptWithPassword(data, pwd)
    if err != nil {
        t.Fatalf("encrypt error: %v", err)
    }

    pt, err := DecryptWithPassword(ct, pwd)
    if err != nil {
        t.Fatalf("decrypt error: %v", err)
    }

    if !bytes.Equal(pt, data) {
        t.Fatalf("plaintext mismatch: got %q", string(pt))
    }
}

func TestEncryptDecryptWithRandomKey_RoundTrip(t *testing.T) {
    data := []byte{0, 1, 2, 3, 4, 5, 250, 251, 252}

    ct, keyB64, err := EncryptWithRandomKey(data)
    if err != nil {
        t.Fatalf("encrypt error: %v", err)
    }

    // Ensure key is valid base64
    if _, err := base64.StdEncoding.DecodeString(keyB64); err != nil {
        t.Fatalf("invalid base64 key: %v", err)
    }

    pt, err := DecryptWithRandomKey(ct, keyB64)
    if err != nil {
        t.Fatalf("decrypt error: %v", err)
    }

    if !bytes.Equal(pt, data) {
        t.Fatalf("plaintext mismatch")
    }
}

func TestDecryptWithRandomKey_InvalidBase64(t *testing.T) {
    data := []byte{1, 2, 3}
    ct, keyB64, err := EncryptWithRandomKey(data)
    if err != nil {
        t.Fatalf("encrypt error: %v", err)
    }
    if len(ct) == 0 || keyB64 == "" {
        t.Fatalf("unexpected empty outputs")
    }
    // Corrupt the key string so base64 decoding fails
    if _, err := DecryptWithRandomKey(ct, "@@not-base64@@"); err == nil {
        t.Fatalf("expected error for invalid base64 key")
    }
}

func TestDecryptWithPassword_AuthFailure(t *testing.T) {
    pwd := "pw"
    msg := []byte("auth will fail after bitflip")
    ct, err := EncryptWithPassword(msg, pwd)
    if err != nil {
        t.Fatalf("encrypt: %v", err)
    }
    // Flip a bit in ciphertext (but keep IV length intact) to trigger GCM auth error
    if len(ct) < 16 {
        t.Fatalf("ciphertext too short for test")
    }
    ctBad := make([]byte, len(ct))
    copy(ctBad, ct)
    ctBad[len(ctBad)-1] ^= 0xFF
    if _, err := DecryptWithPassword(ctBad, pwd); err == nil {
        t.Fatalf("expected auth error on tampered ciphertext")
    }
}

func TestDecryptWithRandomKey_WrongKeyAuthFailure(t *testing.T) {
    data := []byte("plaintext")
    ct, keyB64, err := EncryptWithRandomKey(data)
    if err != nil {
        t.Fatalf("encrypt error: %v", err)
    }
    // Decode, flip a bit, re-encode
    kb, err := base64.StdEncoding.DecodeString(keyB64)
    if err != nil {
        t.Fatalf("decode key: %v", err)
    }
    kb[len(kb)-1] ^= 0x01
    wrong := base64.StdEncoding.EncodeToString(kb)
    if _, err := DecryptWithRandomKey(ct, wrong); err == nil {
        t.Fatalf("expected auth error with wrong key")
    }
}

// failingReader implements io.Reader and always returns an error
type failingReader struct{}

func (f failingReader) Read(p []byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestEncryptWithPassword_IVReadPanic(t *testing.T) {
    // Swap out randReader temporarily to force io.ReadFull error
    old := randReader
    randReader = failingReader{}
    defer func() { randReader = old }()

    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("expected panic when IV generation fails")
        }
    }()
    // This should panic due to failingReader
    _, _ = EncryptWithPassword([]byte("data"), "pwd")
}

func TestDecryptWithPassword_TooShortCiphertextPanics(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("expected panic on too-short ciphertext slice for nonce")
        }
    }()
    _, _ = DecryptWithPassword([]byte{1, 2, 3}, "pwd")
}

func TestEncryptWithRandomKey_KeyGenError(t *testing.T) {
    // Force key generation to fail by making randReader return an error
    old := randReader
    randReader = failingReader{}
    defer func() { randReader = old }()

    if _, _, err := EncryptWithRandomKey([]byte("x")); err == nil {
        t.Fatalf("expected error when random key generation fails")
    }
}
