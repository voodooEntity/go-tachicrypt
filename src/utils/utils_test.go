package utils

import (
    "encoding/hex"
    "strings"
    "testing"
    "bytes"
    "os"
    "io"
)

// captureStdout captures stdout while fn executes and returns the printed string.
func captureStdout(t *testing.T, fn func()) string {
    t.Helper()
    old := os.Stdout
    r, w, err := os.Pipe()
    if err != nil {
        t.Fatalf("pipe: %v", err)
    }
    os.Stdout = w
    defer func() { os.Stdout = old }()

    fn()

    _ = w.Close()
    var buf bytes.Buffer
    if _, err := io.Copy(&buf, r); err != nil {
        t.Fatalf("copy: %v", err)
    }
    _ = r.Close()
    return buf.String()
}

func TestGenerateRandomFilename_HexAndLength(t *testing.T) {
    name1, err := GenerateRandomFilename()
    if err != nil {
        t.Fatalf("GenerateRandomFilename error: %v", err)
    }
    name2, err := GenerateRandomFilename()
    if err != nil {
        t.Fatalf("GenerateRandomFilename error: %v", err)
    }

    if len(name1) != 32 {
        t.Fatalf("expected hex length 32, got %d (%q)", len(name1), name1)
    }
    if _, err := hex.DecodeString(name1); err != nil {
        t.Fatalf("not valid hex: %v", err)
    }
    if name1 == name2 {
        t.Fatalf("expected different filenames across invocations")
    }
}

func TestConcatByteSlices_BasicAndEmpty(t *testing.T) {
    out := ConcatByteSlices([][]byte{[]byte("hello"), []byte(" "), []byte("world")})
    if string(out) != "hello world" {
        t.Fatalf("unexpected concat result: %q", string(out))
    }

    // Empty input
    out2 := ConcatByteSlices(nil)
    if len(out2) != 0 {
        t.Fatalf("expected empty result for nil input")
    }
    out3 := ConcatByteSlices([][]byte{})
    if len(out3) != 0 {
        t.Fatalf("expected empty result for empty slice input")
    }
}

func TestGenerateRandomBytes_LengthRange(t *testing.T) {
    // min <= max
    b, err := GenerateRandomBytes(8, 16)
    if err != nil {
        t.Fatalf("GenerateRandomBytes error: %v", err)
    }
    if l := len(b); l < 8 || l > 16 {
        t.Fatalf("length %d not in [8,16]", l)
    }

    // min > max should be handled by swapping
    b2, err := GenerateRandomBytes(32, 17)
    if err != nil {
        t.Fatalf("GenerateRandomBytes swap error: %v", err)
    }
    if l := len(b2); l < 17 || l > 32 {
        t.Fatalf("length %d not in [17,32] after swap", l)
    }

    // min == max yields fixed length
    b3, err := GenerateRandomBytes(5, 5)
    if err != nil {
        t.Fatalf("GenerateRandomBytes fixed error: %v", err)
    }
    if len(b3) != 5 {
        t.Fatalf("expected length 5, got %d", len(b3))
    }
}

func TestPrintApplicationHeader_OutputsCoreInfo(t *testing.T) {
    version := "v9.9.9-test"
    out := captureStdout(t, func() {
        PrintApplicationHeader(version)
    })
    // Check a few stable substrings (ignore ANSI codes added by prettywriter)
    if !strings.Contains(out, "> Version: "+version) {
        t.Fatalf("header should contain version, got: %q", out)
    }
    if !strings.Contains(out, "> Github: https://github.com/voodooEntity/go-tachicrypt") {
        t.Fatalf("header should contain github link")
    }
    if !strings.Contains(out, "> Author: voodooEntity") {
        t.Fatalf("header should contain author")
    }
}
