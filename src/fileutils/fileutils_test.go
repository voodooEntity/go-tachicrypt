package fileutils

import (
    "bytes"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func TestWriteToFile_NewFileWritesContent(t *testing.T) {
    tmp := t.TempDir()
    path := filepath.Join(tmp, "sample.bin")
    data := []byte{0x00, 0x01, 0x02, 0xFF}

    if err := WriteToFile(path, data); err != nil {
        t.Fatalf("WriteToFile error: %v", err)
    }

    got, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("ReadFile error: %v", err)
    }
    if !bytes.Equal(got, data) {
        t.Fatalf("content mismatch: got %v want %v", got, data)
    }
}

func TestWriteToFile_ErrorOnDirectory(t *testing.T) {
    tmp := t.TempDir()
    // Attempt to write where a directory exists should error
    dirPath := filepath.Join(tmp, "adir")
    if err := os.MkdirAll(dirPath, 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := WriteToFile(dirPath, []byte("x")); err == nil {
        t.Fatalf("expected error when writing to directory path, got nil")
    }
}

func TestObfuscateFileTimestamps_WithinBounds(t *testing.T) {
    tmp := t.TempDir()
    f1 := filepath.Join(tmp, "f1.txt")
    f2 := filepath.Join(tmp, "nested", "f2.txt")
    if err := os.MkdirAll(filepath.Dir(f2), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := os.WriteFile(f1, []byte("a"), 0o644); err != nil {
        t.Fatalf("write f1: %v", err)
    }
    if err := os.WriteFile(f2, []byte("b"), 0o644); err != nil {
        t.Fatalf("write f2: %v", err)
    }

    stat1Before, _ := os.Stat(f1)
    stat2Before, _ := os.Stat(f2)

    // Only files directly under tmp are processed by ObfuscateFileTimestamps
    if err := ObfuscateFileTimestamps(tmp); err != nil {
        t.Fatalf("ObfuscateFileTimestamps error: %v", err)
    }

    stat1After, _ := os.Stat(f1)
    // f2 is in a subdir and won't be touched by current implementation
    stat2After, _ := os.Stat(f2)

    if stat1After.ModTime().Equal(stat1Before.ModTime()) {
        t.Fatalf("expected f1 mtime to change")
    }

    // The random offset is within +/- ~5.8 days; allow a 7-day bound
    now := time.Now()
    lower := now.Add(-7 * 24 * time.Hour)
    upper := now.Add(7 * 24 * time.Hour)
    if mt := stat1After.ModTime(); mt.Before(lower) || mt.After(upper) {
        t.Fatalf("f1 ModTime out of expected bounds: %v not in [%v, %v]", mt, lower, upper)
    }

    // Ensure untouched nested file stays the same (documents current behavior)
    if !stat2After.ModTime().Equal(stat2Before.ModTime()) {
        t.Fatalf("expected f2 mtime to remain unchanged for nested file in current implementation")
    }
}

func TestObfuscateFileTimestamps_NonexistentDir(t *testing.T) {
    if err := ObfuscateFileTimestamps(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
        t.Fatalf("expected error for nonexistent directory")
    }
}
