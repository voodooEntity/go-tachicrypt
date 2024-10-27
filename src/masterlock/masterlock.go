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
	Parts        []PartInfo `json:"parts"`
	FrontPadding int        `json:"front_padding"`
	BackPadding  int        `json:"padding"`
}

func CreateMasterLock(parts []PartInfo, frontPadding int, backPadding int) ([]byte, error) {
	masterLock := MasterLock{
		Parts:        parts,
		FrontPadding: frontPadding,
		BackPadding:  backPadding,
	}

	data, err := json.Marshal(masterLock)
	if err != nil {
		return nil, fmt.Errorf("Could not json encode master lock: %w", err)
	}

	return data, nil
}
