package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "testing"
)

// Test the error branch in main() when hideFunc returns an error
func TestMain_Hide_ErrorPath_Stubbed(t *testing.T) {
    // Prepare temp layout and valid args to pass validation
    tmp := t.TempDir()
    src := filepath.Join(tmp, "file.txt")
    if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
        t.Fatalf("write src: %v", err)
    }
    enc := filepath.Join(tmp, "enc")
    if err := os.MkdirAll(enc, 0o755); err != nil {
        t.Fatalf("mk enc: %v", err)
    }
    os.Setenv("TACHICRYPT_PASSWORD", "p")
    t.Cleanup(func() { os.Unsetenv("TACHICRYPT_PASSWORD") })

    // Stub hooks
    oldHide := hideFunc
    oldExit := exitErrorFn
    called := false
    hideFunc = func(dataPath string, partCount int, outputDir string, prefilledPassword string) error {
        return fmt.Errorf("boom-hide")
    }
    exitErrorFn = func(message string) {
        called = true
        panic("test-exit")
    }
    t.Cleanup(func() { hideFunc = oldHide; exitErrorFn = oldExit })

    // Run main() with valid hide flags
    defer func() { _ = recover() }()
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--hide", "--parts", "2", "--data", src, "--output", enc}
    main()

    if !called {
        t.Fatalf("expected exitErrorFn to be called on hide error")
    }
}

// Test the error branch in main() when unhideFunc returns an error
func TestMain_Unhide_ErrorPath_Stubbed(t *testing.T) {
    tmp := t.TempDir()
    enc := filepath.Join(tmp, "enc")
    out := filepath.Join(tmp, "out")
    if err := os.MkdirAll(enc, 0o755); err != nil { t.Fatalf("mk enc: %v", err) }
    if err := os.MkdirAll(out, 0o755); err != nil { t.Fatalf("mk out: %v", err) }
    os.Setenv("TACHICRYPT_PASSWORD", "p")
    t.Cleanup(func() { os.Unsetenv("TACHICRYPT_PASSWORD") })

    oldUnhide := unhideFunc
    oldExit := exitErrorFn
    called := false
    unhideFunc = func(dataPath string, outputDir string, prefilledPassword string) error {
        return fmt.Errorf("boom-unhide")
    }
    exitErrorFn = func(message string) {
        called = true
        panic("test-exit")
    }
    t.Cleanup(func() { unhideFunc = oldUnhide; exitErrorFn = oldExit })

    defer func() { _ = recover() }()
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--unhide", "--data", enc, "--output", out}
    main()

    if !called {
        t.Fatalf("expected exitErrorFn to be called on unhide error")
    }
}

// Cover the branch in main() where validateFlags fails (missing --parts on hide)
func TestMain_ValidateFails_MissingParts_NoPanic(t *testing.T) {
    oldExit := exitErrorFn
    called := false
    // Do not panic here; allow main() to hit the "return" after validateFlags
    exitErrorFn = func(message string) { called = true }
    defer func() { exitErrorFn = oldExit }()

    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--hide", "--data", "/tmp/x", "--output", "/tmp/y"} // no --parts
    main()
    if !called {
        t.Fatalf("expected exitErrorFn to be called for missing --parts validation failure")
    }
}

// Cover the branch in main() where both --hide and --unhide are set
func TestMain_ValidateFails_BothFlags_NoPanic(t *testing.T) {
    oldExit := exitErrorFn
    called := false
    exitErrorFn = func(message string) { called = true }
    defer func() { exitErrorFn = oldExit }()

    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
    os.Args = []string{"tachicrypt", "--hide", "--unhide", "--parts", "2", "--data", "/tmp/x", "--output", "/tmp/y"}
    main()
    if !called {
        t.Fatalf("expected exitErrorFn to be called when both --hide and --unhide are provided")
    }
}
