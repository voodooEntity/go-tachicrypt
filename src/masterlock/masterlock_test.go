package masterlock

import (
    "encoding/json"
    "testing"
)

func TestCreateMasterLock_JSONRoundTrip(t *testing.T) {
    parts := []PartInfo{
        {Index: 0, Filename: "fileA", Key: "keyA"},
        {Index: 1, Filename: "fileB", Key: "keyB"},
        {Index: 2, Filename: "fileC", Key: "keyC"},
    }
    front := 123
    back := 7

    data, err := CreateMasterLock(parts, front, back)
    if err != nil {
        t.Fatalf("CreateMasterLock error: %v", err)
    }

    // Ensure it is valid JSON and fields map correctly
    var ml MasterLock
    if err := json.Unmarshal(data, &ml); err != nil {
        t.Fatalf("unmarshal error: %v", err)
    }

    if len(ml.Parts) != len(parts) {
        t.Fatalf("parts length mismatch: got %d want %d", len(ml.Parts), len(parts))
    }
    for i := range parts {
        if ml.Parts[i].Index != parts[i].Index || ml.Parts[i].Filename != parts[i].Filename || ml.Parts[i].Key != parts[i].Key {
            t.Fatalf("part %d mismatch: got %+v want %+v", i, ml.Parts[i], parts[i])
        }
    }
    if ml.FrontPadding != front {
        t.Fatalf("FrontPadding mismatch: got %d want %d", ml.FrontPadding, front)
    }
    if ml.BackPadding != back { // JSON tag is "padding"
        t.Fatalf("BackPadding mismatch: got %d want %d", ml.BackPadding, back)
    }
}

func TestCreateMasterLock_EmptyParts(t *testing.T) {
    data, err := CreateMasterLock(nil, 0, 0)
    if err != nil {
        t.Fatalf("CreateMasterLock error: %v", err)
    }

    var ml MasterLock
    if err := json.Unmarshal(data, &ml); err != nil {
        t.Fatalf("unmarshal error: %v", err)
    }

    if len(ml.Parts) != 0 {
        t.Fatalf("expected 0 parts, got %d", len(ml.Parts))
    }
    if ml.FrontPadding != 0 || ml.BackPadding != 0 {
        t.Fatalf("expected zero paddings, got front=%d back=%d", ml.FrontPadding, ml.BackPadding)
    }
}
