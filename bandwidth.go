package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type bandwidthStorage struct {
	MTime  time.Time
	Rx, Tx int
}

func bandwidth() error {
	iface, err := getDefaultIface()
	if err != nil {
		return fmt.Errorf("unable to obtain default iface: %w", err)
	}

	var (
		storage     bandwidthStorage
		storagePath = getStoragePath("bandwidth")
	)
	if err := loadStorage(storagePath, &storage); err != nil {
		return fmt.Errorf("unable to load storage: %w", err)
	}

	rx, err := getRxBytes(iface)
	if err != nil {
		return fmt.Errorf("unable to get received bytes count: %w", err)
	}

	tx, err := getTxBytes(iface)
	if err != nil {
		return fmt.Errorf("unable to get transferred bytes count: %w", err)
	}

	dRx := humanize(rx - storage.Rx)
	dTx := humanize(tx - storage.Tx)
	fmt.Printf("↓↑ %s/s %s/s\n", dRx, dTx)

	storage.Rx = rx
	storage.Tx = tx
	storage.MTime = time.Now()

	if err := saveStorage(storagePath, storage); err != nil {
		return fmt.Errorf("unable to save storage: %w", err)
	}

	return nil
}

func humanize(bytes int) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "ZiB"}
	for i := len(units) - 1; i >= 0; i-- {
		d := 1 << (10 * i)
		if bytes >= d {
			return fmt.Sprintf("%.0f %s", float64(bytes)/float64(d), units[i])
		}
	}

	return "0 B"
}

func getDefaultIface() (string, error) {
	out, err := exec.Command("ip", "route").Output()
	if err != nil {
		return "", fmt.Errorf("unable to run ip: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if !strings.HasPrefix(line, "default via") {
			continue
		}

		tokens := strings.Split(line, " ")
		return tokens[len(tokens)-1], nil
	}

	return string(out), nil
}

func getRxBytes(iface string) (int, error) {
	return getXxBytes(fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", iface))
}
func getTxBytes(iface string) (int, error) {
	return getXxBytes(fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", iface))
}

func getXxBytes(path string) (int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("unable to read net stats file: %w", err)
	}

	ret, err := strconv.Atoi(strings.Trim(string(raw), "\n"))
	if err != nil {
		return 0, fmt.Errorf("unable to parse net stats file: %w", err)
	}

	return ret, nil
}
