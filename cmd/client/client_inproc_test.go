package main

import (
    "flag"
    "os"
    "path/filepath"
    "testing"
)

// TestMain_InProcess_RunSuccessPaths executes hide and unhide within the same
// process by resetting the global FlagSet between invocations. This ensures we
// cover the success branches inside main() without relying solely on subprocess
// runs (which do not attribute coverage to the parent process).
func TestMain_InProcess_RunSuccessPaths(t *testing.T) {
    // Prepare temp filesystem layout
    tmp := t.TempDir()
    src := filepath.Join(tmp, "file.txt")
    if err := os.WriteFile(src, []byte("ok"), 0o644); err != nil {
        t.Fatalf("write src: %v", err)
    }
    enc := filepath.Join(tmp, "enc")
    out := filepath.Join(tmp, "out")
    if err := os.MkdirAll(enc, 0o755); err != nil { t.Fatalf("mk enc: %v", err) }
    if err := os.MkdirAll(out, 0o755); err != nil { t.Fatalf("mk out: %v", err) }

    // 1) Run hide path
    os.Setenv("TACHICRYPT_PASSWORD", "cover-pass")
    t.Cleanup(func() { os.Unsetenv("TACHICRYPT_PASSWORD") })

    // Reset the default FlagSet to avoid conflicts across multiple main() invocations
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--hide", "--parts", "2", "--data", src, "--output", enc}
    main()

    // Ensure artifacts exist (masterlock should be present)
    if _, err := os.Stat(filepath.Join(enc, "masterlock")); err != nil {
        t.Fatalf("expected masterlock after hide: %v", err)
    }

    // 2) Run unhide path
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--unhide", "--data", enc, "--output", out}
    main()

    // Verify restored content
    b, err := os.ReadFile(filepath.Join(out, filepath.Base(src)))
    if err != nil {
        t.Fatalf("read restored: %v", err)
    }
    if string(b) != "ok" {
        t.Fatalf("restored file mismatch: %q", string(b))
    }
}
