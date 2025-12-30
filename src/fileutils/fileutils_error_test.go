package fileutils

import (
    "os"
    "path/filepath"
    "runtime"
    "testing"
)

// TestWriteToFile_WriteError attempts to write to /dev/full (on Unix),
// which accepts opens but fails writes with ENOSPC. If /dev/full does
// not exist (e.g., on non-Unix), the test is skipped.
func TestWriteToFile_WriteError(t *testing.T) {
    if runtime.GOOS == "windows" {
        t.Skip("/dev/full not available on Windows")
    }
    if _, err := os.Stat("/dev/full"); err != nil {
        t.Skip("/dev/full not present; skipping write-error test")
    }
    if err := WriteToFile("/dev/full", []byte("data")); err == nil {
        t.Fatalf("expected error writing to /dev/full, got nil")
    }
}

// TestObfuscateFileTimestamps_TouchError creates a dangling symlink in the
// directory so that os.Chtimes fails when ObfuscateFileTimestamps iterates it.
// If symlink creation is not supported, the test is skipped.
func TestObfuscateFileTimestamps_TouchError(t *testing.T) {
    // Symlinks generally unsupported or require privilege on Windows
    if runtime.GOOS == "windows" {
        t.Skip("symlink behavior differs on Windows; skip")
    }

    dir := t.TempDir()
    // Create a valid file so we also have at least one success candidate
    if err := os.WriteFile(filepath.Join(dir, "ok.txt"), []byte("x"), 0o644); err != nil {
        t.Fatalf("write ok: %v", err)
    }
    // Create a dangling symlink "broken" -> non-existent target
    brokenTarget := filepath.Join(dir, "does-not-exist.txt")
    brokenLink := filepath.Join(dir, "broken")
    if err := os.Symlink(brokenTarget, brokenLink); err != nil {
        t.Skipf("symlink not supported: %v", err)
    }

    if err := ObfuscateFileTimestamps(dir); err == nil {
        t.Fatalf("expected error due to dangling symlink causing Chtimes failure")
    }
}
