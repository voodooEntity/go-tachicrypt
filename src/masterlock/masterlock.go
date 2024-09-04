package masterlock

import (
	"encoding/json"
	"fmt"
)

type PartInfo struct {
	Index    int    `json:"index"`
	Filename string `json:"filename"`
	Key      string `json:"key"`
}

type MasterLock struct {
	Parts   []PartInfo `json:"parts"`
	Padding int        `json:"padding"`
}

func CreateMasterLock(parts []PartInfo, padding int) ([]byte, error) {
	masterLock := MasterLock{
		Parts:   parts,
		Padding: padding,
	}

	data, err := json.Marshal(masterLock)
	if err != nil {
		return nil, fmt.Errorf("Could not json encode master lock: %w", err)
	}

	return data, nil
}
