package main

import "testing"

func TestValidateFlags_MissingParts(t *testing.T) {
    old := exitErrorFn
    called := false
    exitErrorFn = func(string) { called = true }
    t.Cleanup(func() { exitErrorFn = old })

    ok := validateFlags(true, false, -1, "/x", "/y")
    if ok || !called {
        t.Fatalf("expected validation to fail and call exitErrorFn")
    }
}

func TestValidateFlags_BothHideAndUnhide(t *testing.T) {
    old := exitErrorFn
    called := false
    exitErrorFn = func(string) { called = true }
    t.Cleanup(func() { exitErrorFn = old })

    ok := validateFlags(true, true, 3, "/x", "/y")
    if ok || !called {
        t.Fatalf("expected validation to fail for both flags and call exitErrorFn")
    }
}

func TestValidateFlags_MissingDataOrOutput(t *testing.T) {
    old := exitErrorFn
    called := false
    exitErrorFn = func(string) { called = true }
    t.Cleanup(func() { exitErrorFn = old })

    ok := validateFlags(true, false, 2, "", "/y")
    if ok || !called {
        t.Fatalf("expected validation to fail due to missing data/output and call exitErrorFn")
    }
}

func TestValidateFlags_Success(t *testing.T) {
    // Should not call exitErrorFn and return true
    old := exitErrorFn
    called := false
    exitErrorFn = func(string) { called = true }
    t.Cleanup(func() { exitErrorFn = old })

    ok := validateFlags(true, false, 2, "/x", "/y")
    if !ok || called {
        t.Fatalf("expected validation to succeed without calling exitErrorFn")
    }
}
