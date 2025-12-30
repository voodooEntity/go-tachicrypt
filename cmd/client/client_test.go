package main

import (
    "bytes"
    "io"
    "os"
    "strings"
    "testing"
    "flag"
)

// captureStdout captures stdout output while fn runs and returns the printed string.
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

func TestPrintUsage_ContainsCoreSections(t *testing.T) {
    out := captureStdout(t, func() { printUsage() })

    // Strip ANSI escape codes to make assertions robust
    // Basic removal: drop sequences starting with ESC [ and ending with 'm'
    cleaned := out
    cleaned = strings.ReplaceAll(cleaned, "\u001b", "\x1b") // normalize possible encodings
    for {
        i := strings.Index(cleaned, "\x1b[")
        if i < 0 {
            break
        }
        j := strings.Index(cleaned[i:], "m")
        if j < 0 {
            break
        }
        cleaned = cleaned[:i] + cleaned[i+j+1:]
    }

    mustContain := []string{
        "Usage: tachicrypt",
        "Options:",
        "--hide",
        "--unhide",
        "--data",
        "--parts",
        "--output",
        "--help",
        "Examples:",
        "Encrypt data: tachicrypt --hide",
        "Decrypt data: tachicrypt --data",
    }
    for _, s := range mustContain {
        if !strings.Contains(cleaned, s) {
            t.Fatalf("usage output missing %q. Output: %q", s, cleaned)
        }
    }
}

func TestMain_HelpFlagPath(t *testing.T) {
    out := captureStdout(t, func() {
        flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
        os.Args = []string{"tachicrypt", "--help"}
        main()
    })
    if !strings.Contains(out, "Usage: tachicrypt") {
        t.Fatalf("expected usage output when --help is provided; got: %q", out)
    }
}
