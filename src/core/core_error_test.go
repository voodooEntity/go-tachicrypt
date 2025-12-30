package core

import (
    "os"
    "path/filepath"
    "testing"
)

func TestCore_Hide_InvalidPathReturnsError(t *testing.T) {
    tmp := t.TempDir()
    c := New()
    err := c.Hide(filepath.Join(tmp, "definitely-does-not-exist"), 3, tmp, "pass")
    if err == nil {
        t.Fatalf("expected error when hiding nonexistent path")
    }
}

func TestCore_Unhide_WrongPassword(t *testing.T) {
    tmp := t.TempDir()
    // prepare a small file to encrypt
    src := filepath.Join(tmp, "file.txt")
    if err := os.WriteFile(src, []byte("secret"), 0o644); err != nil {
        t.Fatalf("write src: %v", err)
    }

    enc := filepath.Join(tmp, "enc")
    out := filepath.Join(tmp, "out")
    if err := os.MkdirAll(enc, 0o755); err != nil { t.Fatalf("mkdir enc: %v", err) }
    if err := os.MkdirAll(out, 0o755); err != nil { t.Fatalf("mkdir out: %v", err) }

    c := New()
    const right = "right-pass"
    if err := c.Hide(src, 2, enc, right); err != nil {
        t.Fatalf("hide: %v", err)
    }

    // Try to unhide with the wrong password
    if err := c.Unhide(enc, out, "wrong-pass"); err == nil {
        t.Fatalf("expected error when decrypting masterlock with wrong password")
    }
}

func TestCore_Hide_OutputIsAFile_ShouldError(t *testing.T) {
    tmp := t.TempDir()
    // Create a file and attempt to use it as outputDir
    outFile := filepath.Join(tmp, "out-file")
    if err := os.WriteFile(outFile, []byte("x"), 0o644); err != nil {
        t.Fatalf("write outFile: %v", err)
    }
    src := filepath.Join(tmp, "in.txt")
    if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
        t.Fatalf("write src: %v", err)
    }
    c := New()
    if err := c.Hide(src, 2, outFile, "p"); err == nil {
        t.Fatalf("expected error when output path is a file, got nil")
    }
}

func TestCore_Unhide_MissingPartFile_ShouldError(t *testing.T) {
    tmp := t.TempDir()
    src := filepath.Join(tmp, "in.txt")
    if err := os.WriteFile(src, []byte("abc"), 0o644); err != nil { t.Fatalf("write src: %v", err) }
    enc := filepath.Join(tmp, "enc")
    out := filepath.Join(tmp, "out")
    if err := os.MkdirAll(enc, 0o755); err != nil { t.Fatalf("mk enc: %v", err) }
    if err := os.MkdirAll(out, 0o755); err != nil { t.Fatalf("mk out: %v", err) }

    c := New()
    if err := c.Hide(src, 2, enc, "pw"); err != nil { t.Fatalf("hide: %v", err) }

    // Remove one random non-masterlock file from enc to simulate missing part
    entries, err := os.ReadDir(enc)
    if err != nil { t.Fatalf("readdir: %v", err) }
    removed := false
    for _, e := range entries {
        if e.Name() == "masterlock" { continue }
        _ = os.Remove(filepath.Join(enc, e.Name()))
        removed = true
        break
    }
    if !removed {
        t.Fatalf("did not find part file to remove")
    }

    if err := c.Unhide(enc, out, "pw"); err == nil {
        t.Fatalf("expected error when a part file is missing")
    }
}
