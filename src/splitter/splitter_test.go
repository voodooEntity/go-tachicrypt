package splitter

import "testing"

func TestSplitBytesWithPadding_Even(t *testing.T) {
    data := make([]byte, 100)
    parts, pad := SplitBytesWithPadding(data, 5)
    if pad != 0 {
        t.Fatalf("expected padding 0, got %d", pad)
    }
    if len(parts) != 5 {
        t.Fatalf("expected 5 parts, got %d", len(parts))
    }
    for i, p := range parts {
        if len(p) != 20 {
            t.Fatalf("part %d length expected 20, got %d", i, len(p))
        }
    }
}

func TestSplitBytesWithPadding_Odd(t *testing.T) {
    data := make([]byte, 103) // 103 / 5 => part=20 remainder=3
    parts, pad := SplitBytesWithPadding(data, 5)
    if pad != 3 {
        t.Fatalf("expected padding 3, got %d", pad)
    }
    if len(parts) != 5 {
        t.Fatalf("expected 5 parts, got %d", len(parts))
    }
    // All parts 0..3 have 20, last has 20+3 then +3 padding = 26
    for i := 0; i < 4; i++ {
        if len(parts[i]) != 20 {
            t.Fatalf("part %d length expected 20, got %d", i, len(parts[i]))
        }
    }
    if len(parts[4]) != 26 {
        t.Fatalf("last part length expected 26, got %d", len(parts[4]))
    }
}
