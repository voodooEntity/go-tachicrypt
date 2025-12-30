package main

import (
    "flag"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

// runClient runs main() in-process with a fresh FlagSet and a stubbed exitErrorFn.
// It returns stdout/stderr output (captured by caller if desired), whether exitError was invoked, and the exit message.
func runClient(t *testing.T, args ...string) (called bool, msg string) {
    t.Helper()
    // Reset flags between runs
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    // Stub exitErrorFn to avoid os.Exit
    oldExit := exitErrorFn
    var calledFlag bool
    var gotMsg string
    exitErrorFn = func(m string) {
        calledFlag = true
        gotMsg = m
        panic("test-exit")
    }
    t.Cleanup(func() { exitErrorFn = oldExit })

    os.Args = append([]string{"tachicrypt"}, args...)
    defer func() {
        _ = recover() // swallow our test-exit panic
    }()
    main()
    return calledFlag, gotMsg
}

// Validation error branches are covered in validate_flags_test.go.

func TestMain_NoFlags_ShowsUsage(t *testing.T) {
    // capture stdout
    out := captureStdout(t, func() {
        flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
        os.Args = []string{"tachicrypt"}
        main()
    })
    if !strings.Contains(out, "Usage: tachicrypt") {
        t.Fatalf("expected usage output when no flags provided; got: %q", out)
    }
}

func TestMain_Smoke_HideAndUnhide(t *testing.T) {
    // Create temp dirs and a sample file
    tmp := t.TempDir()
    src := filepath.Join(tmp, "hello.txt")
    if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil { t.Fatalf("write src: %v", err) }
    enc := filepath.Join(tmp, "enc")
    out := filepath.Join(tmp, "out")
    if err := os.MkdirAll(enc, 0o755); err != nil { t.Fatalf("mk enc: %v", err) }
    if err := os.MkdirAll(out, 0o755); err != nil { t.Fatalf("mk out: %v", err) }

    os.Setenv("TACHICRYPT_PASSWORD", "smoke-pass")
    t.Cleanup(func() { os.Unsetenv("TACHICRYPT_PASSWORD") })

    // Hide
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--hide", "--parts", "2", "--data", src, "--output", enc}
    main()
    // Unhide
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--unhide", "--data", enc, "--output", out}
    main()

    // Verify restored file content
    restored := filepath.Join(out, "hello.txt")
    b, err := os.ReadFile(restored)
    if err != nil { t.Fatalf("read restored: %v", err) }
    if string(b) != "hello" {
        t.Fatalf("restored content mismatch: %q", string(b))
    }
}

// Note: detailed error paths for unhide are covered via unit tests in core; here we cover
// main() success path and validation errors to exercise CLI branching.
