package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func getStoragePath(name string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("unable to obtain config dir: %w", err))
	}

	return filepath.Join(dir, fmt.Sprintf("i3blocks.%s.json", name))
}

func loadStorage(path string, dst interface{}) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read storage: %w", err)
	}

	if err := json.Unmarshal(raw, &dst); err != nil {
		return fmt.Errorf("invalid data found in storage: %w", err)
	}

	return nil
}

func saveStorage(path string, storage interface{}) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("unable to open storage for writing: %w", err)
	}

	enc := json.NewEncoder(f)
	if err := enc.Encode(storage); err != nil {
		_ = f.Close()
		return fmt.Errorf("unable to encode storage: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("unable to close storage after writing: %w", err)
	}

	return nil
}
