package masterlock

import (
    "errors"
    "testing"
)

func TestCreateMasterLock_MarshalError(t *testing.T) {
    // Save and restore hook
    old := jsonMarshalFn
    jsonMarshalFn = func(v interface{}) ([]byte, error) { return nil, errors.New("marshal boom") }
    t.Cleanup(func() { jsonMarshalFn = old })

    parts := []PartInfo{{Index: 0, Filename: "a", Key: "k"}}
    if _, err := CreateMasterLock(parts, 1, 2); err == nil {
        t.Fatalf("expected error from json marshal hook")
    }
}
