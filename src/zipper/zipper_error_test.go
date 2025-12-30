package zipper

import (
    "archive/zip"
    "bytes"
    "io"
    "os"
    "path/filepath"
    "testing"
)

func TestZip_NonexistentPath(t *testing.T) {
    z := New()
    if _, err := z.Zip("/path/does/not/exist/for-tachicrypt-test"); err == nil {
        t.Fatalf("expected error for nonexistent path")
    }
}

func TestExtract_InvalidBytes(t *testing.T) {
    z := New()
    bad := []byte{0x00, 0x01, 0x02, 0x03} // not a valid zip archive
    if err := z.Extract(bad, t.TempDir()); err == nil {
        t.Fatalf("expected error for invalid zip bytes")
    }
}

// buildZip creates a zip archive with the given entries.
// dirs: list of directory names (trailing slash not required)
// files: map of path -> content
func buildZip(t *testing.T, dirs []string, files map[string]string) []byte {
    t.Helper()
    var buf bytes.Buffer
    w := zip.NewWriter(&buf)
    for _, d := range dirs {
        name := d
        if !bytes.HasSuffix([]byte(name), []byte("/")) {
            name += "/"
        }
        // Ensure directory entry has executable bit so later file open works
        hdr := &zip.FileHeader{Name: name, Method: zip.Deflate}
        hdr.SetMode(0o755 | os.ModeDir)
        if _, err := w.CreateHeader(hdr); err != nil {
            t.Fatalf("create dir entry %q: %v", name, err)
        }
    }
    for p, c := range files {
        fh := &zip.FileHeader{Name: p, Method: zip.Deflate}
        fh.SetMode(0o644)
        f, err := w.CreateHeader(fh)
        if err != nil {
            t.Fatalf("create file entry %q: %v", p, err)
        }
        if _, err := io.WriteString(f, c); err != nil {
            t.Fatalf("write file entry %q: %v", p, err)
        }
    }
    if err := w.Close(); err != nil {
        t.Fatalf("close zip writer: %v", err)
    }
    return buf.Bytes()
}

func TestZip_ReadDirPermissionError(t *testing.T) {
    if os.Geteuid() == 0 {
        t.Skip("running as root; permission tests may not fail as expected")
    }
    tmp := t.TempDir()
    // Create a directory with no read/execute permissions so ReadDir fails
    unrd := filepath.Join(tmp, "sealed")
    if err := os.MkdirAll(unrd, 0o000); err != nil {
        t.Fatalf("mkdir sealed: %v", err)
    }
    // Ensure we can clean it up by chmod back at the end
    t.Cleanup(func() { _ = os.Chmod(unrd, 0o755) })

    z := New()
    if _, err := z.Zip(unrd); err == nil {
        t.Fatalf("expected error when zipping unreadable directory")
    }
}

func TestZip_FileOpenError(t *testing.T) {
    if os.Geteuid() == 0 {
        t.Skip("running as root; permission tests may not fail as expected")
    }
    tmp := t.TempDir()
    // Make a readable directory with one unreadable file
    p := filepath.Join(tmp, "dir")
    if err := os.MkdirAll(p, 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    f := filepath.Join(p, "secret.txt")
    if err := os.WriteFile(f, []byte("topsecret"), 0o000); err != nil {
        t.Fatalf("write secret: %v", err)
    }
    // Restore permissions so cleanup succeeds
    t.Cleanup(func() { _ = os.Chmod(f, 0o644) })

    z := New()
    if _, err := z.Zip(p); err == nil {
        t.Fatalf("expected error when zipping with unreadable file")
    }
}

func TestExtract_DirectoryEntry_SuccessAndMkdirError(t *testing.T) {
    z := New()
    // Build zip that has a directory entry and a file under it
    zipBytes := buildZip(t, []string{"adir"}, map[string]string{"adir/file.txt": "hello"})

    // First, success path
    tmp := t.TempDir()
    if err := z.Extract(zipBytes, tmp); err != nil {
        t.Fatalf("extract with dir entry failed: %v", err)
    }
    b, err := os.ReadFile(filepath.Join(tmp, "adir", "file.txt"))
    if err != nil || string(b) != "hello" {
        t.Fatalf("restored content mismatch: %v %q", err, string(b))
    }

    // Now, error path for MkdirAll when a file blocks directory creation
    tmp2 := t.TempDir()
    // Pre-create a file where a directory is expected
    if err := os.WriteFile(filepath.Join(tmp2, "adir"), []byte("x"), 0o644); err != nil {
        t.Fatalf("precreate file for mkdir error: %v", err)
    }
    if err := z.Extract(zipBytes, tmp2); err == nil {
        t.Fatalf("expected error when MkdirAll tries to create over a file")
    }
}

func TestExtract_FileParentMkdirError_AndCreateError(t *testing.T) {
    z := New()
    zipBytes := buildZip(t, nil, map[string]string{"a/b.txt": "data"})

    // Parent mkdir error: create a file at dest/a so MkdirAll(dest/a) fails
    tmp := t.TempDir()
    if err := os.WriteFile(filepath.Join(tmp, "a"), []byte("block"), 0o644); err != nil {
        t.Fatalf("precreate a: %v", err)
    }
    if err := z.Extract(zipBytes, tmp); err == nil {
        t.Fatalf("expected error when parent MkdirAll fails for file")
    }

    // Create error: allow parent mkdir, but make final outFile a directory
    tmp2 := t.TempDir()
    if err := os.MkdirAll(filepath.Join(tmp2, "a"), 0o755); err != nil {
        t.Fatalf("mkdir parent: %v", err)
    }
    // Create directory at outFile path to cause os.Create fail
    if err := os.MkdirAll(filepath.Join(tmp2, "a", "b.txt"), 0o755); err != nil {
        t.Fatalf("mkdir at file path: %v", err)
    }
    if err := z.Extract(zipBytes, tmp2); err == nil {
        t.Fatalf("expected error when os.Create is called on an existing directory")
    }
}

// failing implementations for hooks
type errCloser struct{}
type errReader struct{}
type errCopier struct{}

func (e errCloser) Close() error { return io.ErrUnexpectedEOF }
func (e errReader) Read(p []byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (e errCopier) Read(p []byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestZip_CloseWriterError(t *testing.T) {
    // Build a small file to zip
    tmp := t.TempDir()
    f := filepath.Join(tmp, "x.txt")
    if err := os.WriteFile(f, []byte("x"), 0o644); err != nil { t.Fatalf("write: %v", err) }

    // Force closeZipWriterFn to error
    oldClose := closeZipWriterFn
    closeZipWriterFn = func(w *zip.Writer) error { return io.ErrUnexpectedEOF }
    t.Cleanup(func() { closeZipWriterFn = oldClose })

    z := New()
    if _, err := z.Zip(f); err == nil {
        t.Fatalf("expected error from closeZipWriterFn")
    }
}

func TestZip_IoCopyError(t *testing.T) {
    // Prepare a file to zip
    tmp := t.TempDir()
    f := filepath.Join(tmp, "x.txt")
    if err := os.WriteFile(f, []byte("x"), 0o644); err != nil { t.Fatalf("write: %v", err) }

    // Force ioCopyFn to fail
    oldCopy := ioCopyFn
    ioCopyFn = func(dst io.Writer, src io.Reader) (int64, error) { return 0, io.ErrUnexpectedEOF }
    t.Cleanup(func() { ioCopyFn = oldCopy })

    z := New()
    if _, err := z.Zip(f); err == nil {
        t.Fatalf("expected error from ioCopyFn during zip")
    }
}

func TestExtract_FileOpenErrorAndCopyError(t *testing.T) {
    // Build a simple zip with one file
    data := buildZip(t, nil, map[string]string{"f.txt": "content"})

    // First, error from zipFileOpenFn
    oldOpen := zipFileOpenFn
    zipFileOpenFn = func(f *zip.File) (io.ReadCloser, error) { return nil, io.ErrUnexpectedEOF }
    t.Cleanup(func() { zipFileOpenFn = oldOpen })

    z := New()
    if err := z.Extract(data, t.TempDir()); err == nil {
        t.Fatalf("expected error from zipFileOpenFn")
    }

    // Next, let open succeed but fail during copy
    zipFileOpenFn = func(f *zip.File) (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader([]byte("abc"))), nil }
    oldCopy := ioCopyFn
    ioCopyFn = func(dst io.Writer, src io.Reader) (int64, error) { return 0, io.ErrUnexpectedEOF }
    t.Cleanup(func() { ioCopyFn = oldCopy })
    if err := z.Extract(data, t.TempDir()); err == nil {
        t.Fatalf("expected error from ioCopyFn during extract")
    }
}

func TestZip_CreateEntryError(t *testing.T) {
    // Prepare a file to zip
    tmp := t.TempDir()
    f := filepath.Join(tmp, "x.txt")
    if err := os.WriteFile(f, []byte("x"), 0o644); err != nil { t.Fatalf("write: %v", err) }

    // Stub createZipEntryFn to return error
    oldCreate := createZipEntryFn
    createZipEntryFn = func(w *zip.Writer, name string) (io.Writer, error) { return nil, io.ErrUnexpectedEOF }
    t.Cleanup(func() { createZipEntryFn = oldCreate })

    z := New()
    if _, err := z.Zip(f); err == nil {
        t.Fatalf("expected error from createZipEntryFn")
    }
}
