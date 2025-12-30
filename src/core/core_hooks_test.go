package core

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "testing"

    ml "github.com/voodooEntity/go-tachicrypt/src/masterlock"
)

// helper to create a tiny input file and dirs
func mkInputEnv(t *testing.T) (src string, encDir string, outDir string) {
    t.Helper()
    tmp := t.TempDir()
    src = filepath.Join(tmp, "in.txt")
    if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
        t.Fatalf("write src: %v", err)
    }
    encDir = filepath.Join(tmp, "enc")
    outDir = filepath.Join(tmp, "out")
    if err := os.MkdirAll(encDir, 0o755); err != nil { t.Fatalf("mk enc: %v", err) }
    if err := os.MkdirAll(outDir, 0o755); err != nil { t.Fatalf("mk out: %v", err) }
    return
}

// restore all hooks to real implementations
// note: we prefer per-test t.Cleanup restorations rather than a global reset

func TestCore_Hide_UsesPromptPasswordBranch(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    // stub prompt to avoid TTY
    oldPrompt := promptPasswordFn
    promptPasswordFn = func(string) string { return "prompt-pass" }
    oldObf := obfuscateFileTimestampsFn
    obfuscateFileTimestampsFn = func(string) error { return nil }
    t.Cleanup(func() { promptPasswordFn = oldPrompt; obfuscateFileTimestampsFn = oldObf })

    c := New()
    if err := c.Hide(src, 2, enc, ""); err != nil {
        t.Fatalf("Hide with prompt branch failed: %v", err)
    }
}

func TestCore_Hide_ErrorFromRandomFrontPadding(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := genRandomBytesFn
    genRandomBytesFn = func(min, max int) ([]byte, error) { return nil, errors.New("rng fail") }
    t.Cleanup(func() { genRandomBytesFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from random padding generation")
    }
}

func TestCore_Hide_ErrorFromEncryptPart(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := encryptWithRandomKeyFn
    encryptWithRandomKeyFn = func([]byte) ([]byte, string, error) { return nil, "", errors.New("enc part") }
    t.Cleanup(func() { encryptWithRandomKeyFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from encrypt part")
    }
}

func TestCore_Hide_ErrorFromGenerateFilename(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := generateRandomFilenameFn
    generateRandomFilenameFn = func() (string, error) { return "", errors.New("name gen") }
    t.Cleanup(func() { generateRandomFilenameFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from filename generation")
    }
}

func TestCore_Hide_ErrorFromWritePart(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := writeToFileFn
    writeToFileFn = func(path string, b []byte) error {
        // Fail for part writes (anything not masterlock)
        if filepath.Base(path) != "masterlock" { return errors.New("write part") }
        return old(path, b)
    }
    t.Cleanup(func() { writeToFileFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from write part")
    }
}

func TestCore_Hide_ErrorFromCreateMasterlock(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := createMasterLockFn
    createMasterLockFn = func(_ []ml.PartInfo, _ int, _ int) ([]byte, error) { return nil, errors.New("mk mlock") }
    t.Cleanup(func() { createMasterLockFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from CreateMasterLock")
    }
}

func TestCore_Hide_ErrorFromEncryptMasterlock(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := encryptWithPasswordFn
    encryptWithPasswordFn = func(data []byte, pwd string) ([]byte, error) { return nil, errors.New("enc mlock") }
    t.Cleanup(func() { encryptWithPasswordFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from EncryptWithPassword on masterlock")
    }
}

func TestCore_Hide_ErrorFromWriteMasterlock(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := writeToFileFn
    writeToFileFn = func(path string, b []byte) error {
        if filepath.Base(path) == "masterlock" { return errors.New("write mlock") }
        return old(path, b)
    }
    t.Cleanup(func() { writeToFileFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error writing masterlock")
    }
}

func TestCore_Hide_ErrorFromObfuscateTimestamps(t *testing.T) {
    src, enc, _ := mkInputEnv(t)
    old := obfuscateFileTimestampsFn
    obfuscateFileTimestampsFn = func(string) error { return errors.New("obf") }
    t.Cleanup(func() { obfuscateFileTimestampsFn = old })
    c := New()
    if err := c.Hide(src, 2, enc, "p"); err == nil {
        t.Fatalf("expected error from obfuscate timestamps")
    }
}

func TestCore_Unhide_ErrorReadingMasterlock(t *testing.T) {
    _, enc, out := mkInputEnv(t)
    // Ensure enc dir exists but without masterlock file; force readFileFn error regardless
    old := readFileFn
    readFileFn = func(string) ([]byte, error) { return nil, errors.New("read mlock") }
    t.Cleanup(func() { readFileFn = old })
    c := New()
    if err := c.Unhide(enc, out, "p"); err == nil {
        t.Fatalf("expected error reading masterlock")
    }
}

func TestCore_Unhide_ErrorDecryptingMasterlock(t *testing.T) {
    _, enc, out := mkInputEnv(t)
    // place a dummy masterlock file (content irrelevant due to hook)
    if err := os.WriteFile(filepath.Join(enc, "masterlock"), []byte{1,2,3}, 0o644); err != nil { t.Fatalf("write mlock: %v", err) }
    old := decryptWithPasswordFn
    decryptWithPasswordFn = func([]byte, string) ([]byte, error) { return nil, errors.New("dec mlock") }
    t.Cleanup(func() { decryptWithPasswordFn = old })
    c := New()
    if err := c.Unhide(enc, out, "p"); err == nil {
        t.Fatalf("expected error decrypting masterlock")
    }
}

func TestCore_Unhide_ErrorJSONUnmarshal(t *testing.T) {
    _, enc, out := mkInputEnv(t)
    if err := os.WriteFile(filepath.Join(enc, "masterlock"), []byte{1}, 0o644); err != nil { t.Fatalf("write mlock: %v", err) }
    oldDec := decryptWithPasswordFn
    decryptWithPasswordFn = func([]byte, string) ([]byte, error) { return []byte("not json"), nil }
    oldUn := jsonUnmarshalFn
    jsonUnmarshalFn = func([]byte, interface{}) error { return errors.New("bad json") }
    t.Cleanup(func() { decryptWithPasswordFn = oldDec; jsonUnmarshalFn = oldUn })
    c := New()
    if err := c.Unhide(enc, out, "p"); err == nil {
        t.Fatalf("expected error from json unmarshal")
    }
}

func TestCore_Unhide_ErrorDecryptingPart(t *testing.T) {
    _, enc, out := mkInputEnv(t)
    // Build a fake masterlock JSON with one part
    mlk := ml.MasterLock{Parts: []ml.PartInfo{{Index:0, Filename:"p1", Key:"k"}}, FrontPadding:0, BackPadding:0}
    data, _ := json.Marshal(mlk)
    if err := os.WriteFile(filepath.Join(enc, "masterlock"), []byte{1}, 0o644); err != nil { t.Fatalf("write mlock: %v", err) }
    // Hook to return our JSON regardless of file contents
    oldDec := decryptWithPasswordFn
    decryptWithPasswordFn = func([]byte, string) ([]byte, error) { return data, nil }
    t.Cleanup(func() { decryptWithPasswordFn = oldDec })
    // Create part file p1
    if err := os.WriteFile(filepath.Join(enc, "p1"), []byte{9,9,9}, 0o644); err != nil { t.Fatalf("write part: %v", err) }
    // Hook decrypt part to error
    oldPart := decryptWithRandomKeyFn
    decryptWithRandomKeyFn = func([]byte, string) ([]byte, error) { return nil, errors.New("dec part") }
    t.Cleanup(func() { decryptWithRandomKeyFn = oldPart })
    c := New()
    if err := c.Unhide(enc, out, "p"); err == nil {
        t.Fatalf("expected error decrypting a part")
    }
}

func TestCore_Unhide_UsesPromptPasswordBranch(t *testing.T) {
    src, enc, out := mkInputEnv(t)
    // First, create encrypted set with known password using normal flow
    c := New()
    if err := c.Hide(src, 2, enc, "known-pass"); err != nil {
        t.Fatalf("hide failed: %v", err)
    }
    // Now unhide but force prompt branch to supply the password
    oldPrompt := promptPasswordFn
    promptPasswordFn = func(string) string { return "known-pass" }
    t.Cleanup(func() { promptPasswordFn = oldPrompt })
    if err := c.Unhide(enc, out, ""); err != nil {
        t.Fatalf("unhide with prompt branch failed: %v", err)
    }
    // Verify round-trip
    restored := filepath.Join(out, filepath.Base(src))
    b, err := os.ReadFile(restored)
    if err != nil { t.Fatalf("read restored: %v", err) }
    if string(b) != "hello" {
        t.Fatalf("restored content mismatch: %q", string(b))
    }
}

func TestCore_Unhide_ErrorOnExtract(t *testing.T) {
    _, enc, out := mkInputEnv(t)
    // Make a masterlock with one part
    mlk := ml.MasterLock{Parts: []ml.PartInfo{{Index:0, Filename:"p1", Key:"k"}}, FrontPadding:0, BackPadding:0}
    data, _ := json.Marshal(mlk)
    if err := os.WriteFile(filepath.Join(enc, "masterlock"), []byte{1}, 0o644); err != nil { t.Fatalf("write mlock: %v", err) }
    // Decrypt masterlock returns our JSON
    oldDec := decryptWithPasswordFn
    decryptWithPasswordFn = func([]byte, string) ([]byte, error) { return data, nil }
    t.Cleanup(func() { decryptWithPasswordFn = oldDec })
    // Create part file p1 and make part decryption return garbage (invalid zip)
    if err := os.WriteFile(filepath.Join(enc, "p1"), []byte{9}, 0o644); err != nil { t.Fatalf("write part: %v", err) }
    oldPart := decryptWithRandomKeyFn
    decryptWithRandomKeyFn = func([]byte, string) ([]byte, error) { return []byte{0x00, 0x01, 0x02}, nil }
    t.Cleanup(func() { decryptWithRandomKeyFn = oldPart })
    c := New()
    if err := c.Unhide(enc, out, "p"); err == nil {
        t.Fatalf("expected error from Extract on invalid zip bytes")
    }
}
