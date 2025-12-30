package zipper

import (
    "bytes"
    "io/fs"
    "os"
    "path/filepath"
    "testing"
)

func writeFile(t *testing.T, path string, data []byte) {
    t.Helper()
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := os.WriteFile(path, data, 0o644); err != nil {
        t.Fatalf("write file: %v", err)
    }
}

func readFile(t *testing.T, path string) []byte {
    t.Helper()
    b, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("read file: %v", err)
    }
    return b
}

func TestZipAndExtract_File(t *testing.T) {
    tmp := t.TempDir()
    srcFile := filepath.Join(tmp, "hello.txt")
    content := []byte("hello zipper file test\nwith multiple lines\n")
    writeFile(t, srcFile, content)

    z := New()
    zipBytes, err := z.Zip(srcFile)
    if err != nil {
        t.Fatalf("zip file error: %v", err)
    }

    dest := filepath.Join(tmp, "out")
    if err := z.Extract(zipBytes, dest); err != nil {
        t.Fatalf("extract error: %v", err)
    }

    restored := filepath.Join(dest, filepath.Base(srcFile))
    got := readFile(t, restored)
    if !bytes.Equal(got, content) {
        t.Fatalf("restored content mismatch")
    }
}

func TestZipAndExtract_DirectoryTree(t *testing.T) {
    tmp := t.TempDir()
    root := filepath.Join(tmp, "tree")
    // create nested structure
    writeFile(t, filepath.Join(root, "a", "file1.txt"), []byte("alpha"))
    writeFile(t, filepath.Join(root, "a", "b", "file2.txt"), []byte("beta"))
    writeFile(t, filepath.Join(root, "c", "file3.txt"), []byte("gamma"))

    z := New()
    zipBytes, err := z.Zip(root)
    if err != nil {
        t.Fatalf("zip dir error: %v", err)
    }

    dest := filepath.Join(tmp, "out")
    if err := z.Extract(zipBytes, dest); err != nil {
        t.Fatalf("extract error: %v", err)
    }

    restoredRoot := filepath.Join(dest, filepath.Base(root))

    // compare trees
    wantFiles := map[string][]byte{}
    err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        rel, _ := filepath.Rel(root, path)
        wantFiles[rel] = readFile(t, path)
        return nil
    })
    if err != nil {
        t.Fatalf("walk src: %v", err)
    }

    gotFiles := map[string][]byte{}
    err = filepath.WalkDir(restoredRoot, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        rel, _ := filepath.Rel(restoredRoot, path)
        gotFiles[rel] = readFile(t, path)
        return nil
    })
    if err != nil {
        t.Fatalf("walk dst: %v", err)
    }

    if len(gotFiles) != len(wantFiles) {
        t.Fatalf("file count mismatch: got %d want %d", len(gotFiles), len(wantFiles))
    }
    for rel, want := range wantFiles {
        got, ok := gotFiles[rel]
        if !ok {
            t.Fatalf("missing file in restored tree: %s", rel)
        }
        if !bytes.Equal(got, want) {
            t.Fatalf("content mismatch for %s", rel)
        }
    }
}

func TestZipAndExtract_EmptyDirectory(t *testing.T) {
    tmp := t.TempDir()
    empty := filepath.Join(tmp, "emptydir")
    if err := os.MkdirAll(empty, 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    z := New()
    zipBytes, err := z.Zip(empty)
    if err != nil {
        t.Fatalf("zip empty dir error: %v", err)
    }
    dest := filepath.Join(tmp, "out")
    if err := z.Extract(zipBytes, dest); err != nil {
        t.Fatalf("extract empty dir error: %v", err)
    }
    // Current implementation doesn't create entries for empty dirs when zipping;
    // just assert that no files were created and extraction succeeded.
    count := 0
    _ = filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
        if err != nil { return err }
        if !info.IsDir() { count++ }
        return nil
    })
    if count != 0 {
        t.Fatalf("expected no files extracted from empty dir, got %d", count)
    }
}
