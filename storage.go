package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func getStoragePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("unable to obtain config dir: %w", err))
	}

	return filepath.Join(dir, "i3blocks.json")
}

type storage struct {
	MTime  time.Time
	Rx, Tx int
}

func loadStorage(path string) (storage, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return storage{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return storage{}, fmt.Errorf("unable to open storage for reading: %w", err)
	}
	defer f.Close()

	var (
		ret storage
		dec = json.NewDecoder(f)
	)

	if err := dec.Decode(&ret); err != nil {
		return storage{}, fmt.Errorf("invalid data found in storage: %w", err)
	}

	return ret, nil
}

func (s *storage) save(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("unable to open storage for writing: %w", err)
	}

	s.MTime = time.Now()
	enc := json.NewEncoder(f)
	if err := enc.Encode(s); err != nil {
		_ = f.Close()
		return fmt.Errorf("unable to encode storage: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("unable to close storage after writing: %w", err)
	}

	return nil
}
