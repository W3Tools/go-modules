package gmar

import (
	"fmt"
	"os"
)

type ArweaveManifest struct {
	Manifest string                         `json:"manifest"`
	Version  string                         `json:"version"`
	Index    ArweaveManifestIndex           `json:"index"`
	Paths    map[string]ArweaveManifestPath `json:"paths"`
}

type ArweaveManifestIndex struct {
	Path string `json:"path"`
}

type ArweaveManifestPath struct {
	ID string `json:"id"`
}

func NewManifest() *ArweaveManifest {
	return &ArweaveManifest{
		Manifest: "arweave/paths", // Default
		Version:  "0.1.0",         // Default
		Index:    ArweaveManifestIndex{Path: ""},
		Paths:    make(map[string]ArweaveManifestPath),
	}
}

func WriteManifest(name string, data []byte) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile err %v", err)
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("f.Write err %v", err)
	}

	return nil
}
