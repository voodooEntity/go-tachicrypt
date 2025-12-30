package encryptor

import (
    "bytes"
    "encoding/base64"
    "testing"
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
