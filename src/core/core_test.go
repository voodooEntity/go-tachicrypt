package core

import (
    "bytes"
    "io/fs"
    "os"
    "path/filepath"
    "testing"
)

// writeFile is a small helper to create a file with content, ensuring parent dirs exist.
func writeFile(t *testing.T, path string, data []byte) {
    t.Helper()
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := os.WriteFile(path, data, 0o644); err != nil {
        t.Fatalf("write file: %v", err)
    }
}

// collectFiles returns a map of relative path -> content for all files under root.
func collectFiles(t *testing.T, root string) map[string][]byte {
    t.Helper()
    files := map[string][]byte{}
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        rel, _ := filepath.Rel(root, path)
        b, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        files[rel] = b
        return nil
    })
    if err != nil {
        t.Fatalf("walk %s: %v", root, err)
    }
    return files
}

func TestCore_RoundTrip_File(t *testing.T) {
    tmp := t.TempDir()
    // prepare a sample file
    srcFile := filepath.Join(tmp, "hello.txt")
    content := []byte("hello from core roundtrip\nwith multiple lines\n")
    writeFile(t, srcFile, content)

    // encryption output dir
    encDir := filepath.Join(tmp, "enc")
    if err := os.MkdirAll(encDir, 0o755); err != nil {
        t.Fatalf("mkdir enc: %v", err)
    }
    // decryption output dir
    outDir := filepath.Join(tmp, "out")
    if err := os.MkdirAll(outDir, 0o755); err != nil {
        t.Fatalf("mkdir out: %v", err)
    }

    c := New()
    const pwd = "unit-test-pass"
    if err := c.Hide(srcFile, 3, encDir, pwd); err != nil {
        t.Fatalf("Hide error: %v", err)
    }
    if err := c.Unhide(encDir, outDir, pwd); err != nil {
        t.Fatalf("Unhide error: %v", err)
    }

    // Restored file should appear under out/<basename>
    restored := filepath.Join(outDir, filepath.Base(srcFile))
    got, err := os.ReadFile(restored)
    if err != nil {
        t.Fatalf("read restored: %v", err)
    }
    if !bytes.Equal(got, content) {
        t.Fatalf("restored content mismatch")
    }
}

func TestCore_RoundTrip_Directory(t *testing.T) {
    tmp := t.TempDir()
    // prepare a small directory tree
    srcRoot := filepath.Join(tmp, "tree")
    writeFile(t, filepath.Join(srcRoot, "a", "file1.txt"), []byte("alpha"))
    writeFile(t, filepath.Join(srcRoot, "a", "b", "file2.txt"), []byte("beta"))
    writeFile(t, filepath.Join(srcRoot, "c", "file3.txt"), []byte("gamma"))

    encDir := filepath.Join(tmp, "enc")
    if err := os.MkdirAll(encDir, 0o755); err != nil {
        t.Fatalf("mkdir enc: %v", err)
    }
    outDir := filepath.Join(tmp, "out")
    if err := os.MkdirAll(outDir, 0o755); err != nil {
        t.Fatalf("mkdir out: %v", err)
    }

    c := New()
    const pwd = "unit-test-pass"
    if err := c.Hide(srcRoot, 4, encDir, pwd); err != nil {
        t.Fatalf("Hide error: %v", err)
    }
    if err := c.Unhide(encDir, outDir, pwd); err != nil {
        t.Fatalf("Unhide error: %v", err)
    }

    // Zipper extracts under out/<basename(srcRoot)>
    restoredRoot := filepath.Join(outDir, filepath.Base(srcRoot))
    got := collectFiles(t, restoredRoot)
    want := collectFiles(t, srcRoot)

    if len(got) != len(want) {
        t.Fatalf("file count mismatch: got %d want %d", len(got), len(want))
    }
    for rel, wb := range want {
        gb, ok := got[rel]
        if !ok {
            t.Fatalf("missing restored file: %s", rel)
        }
        if !bytes.Equal(gb, wb) {
            t.Fatalf("content mismatch for %s", rel)
        }
    }
}
