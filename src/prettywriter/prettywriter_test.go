package prettywriter

import (
    "bytes"
    "io"
    "os"
    "strings"
    "testing"
)

// captureOutput redirects stdout while fn runs and returns what was printed.
func captureOutput(t *testing.T, fn func()) string {
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

func TestWrapMessage_Basic(t *testing.T) {
    msg := "The quick brown fox jumps over the lazy dog"
    lines := wrapMessage(msg, 10)
    if len(lines) == 0 {
        t.Fatalf("expected some lines")
    }
    for i, ln := range lines {
        if len(ln) > 10 {
            t.Fatalf("line %d exceeds width: %q (len=%d)", i, ln, len(ln))
        }
    }
    // Ensure all words are present (order preserved).
    joined := strings.Join(lines, " ")
    // normalize multiple spaces
    f := func(r rune) bool { return r == ' ' || r == '\n' || r == '\t' }
    wantWords := strings.Fields(msg)
    gotWords := strings.FieldsFunc(joined, f)
    if len(gotWords) != len(wantWords) {
        t.Fatalf("word count mismatch: got %d want %d", len(gotWords), len(wantWords))
    }
    for i := range wantWords {
        if gotWords[i] != wantWords[i] {
            t.Fatalf("word %d mismatch: got %q want %q", i, gotWords[i], wantWords[i])
        }
    }
}

func TestGetBorder_SingleVsDouble(t *testing.T) {
    width := 15
    single := getBorder(width, SingleLine)
    if len(single) != width {
        t.Fatalf("single border len=%d want=%d", len(single), width)
    }
    for _, ch := range single {
        if ch != '-' {
            t.Fatalf("single border contains unexpected char: %q", string(ch))
        }
    }

    dbl := getBorder(width, DoubleLine)
    if len(dbl) != width {
        t.Fatalf("double border len=%d want=%d", len(dbl), width)
    }
    for _, ch := range dbl {
        if ch != '=' {
            t.Fatalf("double border contains unexpected char: %q", string(ch))
        }
    }
}

func TestWriteAndWriteln_Output(t *testing.T) {
    out := captureOutput(t, func() {
        Write("hello", Green, BlackBG)
    })
    if !strings.Contains(out, "hello") {
        t.Fatalf("Write did not output message: %q", out)
    }
    if strings.Contains(out, "\n") {
        t.Fatalf("Write should not include newline")
    }
    // Expect ANSI prefix and reset suffix
    if !strings.Contains(out, "\x1b[") || !strings.Contains(out, "\x1b[0m") {
        t.Fatalf("Write should include ANSI color codes: %q", out)
    }

    out2 := captureOutput(t, func() {
        Writeln("world", Green, BlackBG)
    })
    if !strings.HasSuffix(out2, "\n") {
        t.Fatalf("Writeln should end with newline: %q", out2)
    }
}

func TestWriteInBox_PrintsBordersAndMessage(t *testing.T) {
    width := 20
    msg := "hello world from prettywriter"
    out := captureOutput(t, func() {
        WriteInBox(width, msg, Green, BlackBG, DoubleLine)
    })
    lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
    // remove any trailing empty lines
    for len(lines) > 0 && lines[len(lines)-1] == "" {
        lines = lines[:len(lines)-1]
    }
    if len(lines) < 3 {
        t.Fatalf("expected at least 3 lines, got %d: %q", len(lines), out)
    }
    if !strings.HasPrefix(lines[0], "┌") || !strings.HasSuffix(lines[0], "┐") {
        t.Fatalf("top border malformed: %q", lines[0])
    }
    last := lines[len(lines)-1]
    if !strings.HasPrefix(last, "└") || !strings.HasSuffix(last, "┘") {
        t.Fatalf("bottom border malformed: %q", last)
    }
    if !strings.Contains(out, "hello") {
        t.Fatalf("box output should contain message: %q", out)
    }
}
