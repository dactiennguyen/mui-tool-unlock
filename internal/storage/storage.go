package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"muitoolunlock/internal/types"
)

// LoadUnlockData loads unlock data from local file
func LoadUnlockData() *types.UnlockData {
	baseDir, err := os.Getwd()
	if err != nil {
		return &types.UnlockData{}
	}
	dataFile := filepath.Join(baseDir, "miunlockdata.json")

	data := &types.UnlockData{}
	if fileData, err := os.ReadFile(dataFile); err == nil {
		json.Unmarshal(fileData, data)
	}

	return data
}

// SaveUnlockData saves unlock data to local file
func SaveUnlockData(data *types.UnlockData) {
	baseDir, err := os.Getwd()
	if err != nil {
		return
	}
	dataFile := filepath.Join(baseDir, "miunlockdata.json")

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(dataFile, jsonData, 0644)
}
