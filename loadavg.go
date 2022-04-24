package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type stats struct {
	User, System, Idle, Wait int
}

func (a stats) asPercentage() stats {
	total := a.User + a.System + a.Idle + a.Wait

	return stats{
		User:   (100 * a.User) / total,
		System: (100 * a.System) / total,
		Idle:   (100 * a.Idle) / total,
		Wait:   (100 * a.Wait) / total,
	}
}

func (a stats) sub(b stats) stats {
	return stats{
		User:   a.User - b.User,
		System: a.System - b.System,
		Idle:   a.Idle - b.Idle,
		Wait:   a.Wait - b.Wait,
	}
}

func loadavg() error {
	loads, err := getLoads()
	if err != nil {
		return fmt.Errorf("unable to get loads: %w", err)
	}

	cpus := float64(runtime.NumCPU())
	for _, v := range loads {
		fmt.Printf("%s ", gradientFloat(v, 0, cpus, "%.2f"))
	}

	stats, err := computeStats()
	if err != nil {
		return fmt.Errorf("unable to get CPU stats: %w", err)
	}

	pct := stats.asPercentage()
	fmt.Printf(
		"(%s/%s/%s/%d)\n",
		gradientFloat(float64(pct.User), -1, 100, "%.0f"),
		gradientFloat(float64(pct.System), -1, 100, "%.0f"),
		gradientFloat(float64(pct.Wait), -1, 100, "%.0f"),
		pct.Idle,
	)

	return nil
}

func gradientFloat(v, min, max float64, vFormat string) string {
	format := fmt.Sprintf("<span color='%%s'>%s</span>", vFormat)

	switch {
	case v <= min:
		return fmt.Sprintf(format, "green", v)
	case v >= max:
		fallthrough
	case v >= min+((max-min)*.75):
		return fmt.Sprintf(format, "red", v)
	case v >= min+((max-min)*.50):
		return fmt.Sprintf(format, "orange", v)
	case v >= min+((max-min)*.25):
		return fmt.Sprintf(format, "yellow", v)
	default:
		return fmt.Sprintf(vFormat, v)
	}
}

func computeStats() (stats, error) {
	var (
		storage     stats
		storagePath = getStoragePath("stats")
	)
	if err := loadStorage(storagePath, &storage); err != nil {
		return stats{}, fmt.Errorf("unable to load storage: %w", err)
	}

	curStats, err := loadStatsFromProc()
	if err != nil {
		return stats{}, fmt.Errorf("unable to read stats from proc: %w", err)
	}

	if err := saveStorage(storagePath, curStats); err != nil {
		return stats{}, fmt.Errorf("unable to save storage: %w", err)
	}

	return curStats.sub(storage), nil
}

func loadStatsFromProc() (stats, error) {
	raw, err := os.ReadFile("/proc/stat")
	if err != nil {
		return stats{}, fmt.Errorf("unable to read /proc/stat: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			return parseProcStatCPULine(line)
		}
	}

	return stats{}, fmt.Errorf("unable to find CPU info in /proc/stat: %w", err)
}

func parseProcStatCPULine(line string) (stats, error) {
	//            2                     4       5     6
	// "cpu", "", user, nice (ignored), system, idle, iowait
	tokens := strings.Split(line, " ")
	if len(tokens) < 6 {
		return stats{}, errors.New("not enough enough entries in CPU line")
	}

	var (
		ret stats
		err error
	)

	ret.User, err = strconv.Atoi(tokens[2])
	if err != nil {
		return stats{}, fmt.Errorf("unable to parse user: %w", err)
	}

	ret.System, err = strconv.Atoi(tokens[4])
	if err != nil {
		return stats{}, fmt.Errorf("unable to parse system: %w", err)
	}

	ret.Idle, err = strconv.Atoi(tokens[5])
	if err != nil {
		return stats{}, fmt.Errorf("unable to parse idle: %w", err)
	}

	ret.Wait, err = strconv.Atoi(tokens[6])
	if err != nil {
		return stats{}, fmt.Errorf("unable to parse iowait: %w", err)
	}

	return ret, nil
}

func getLoads() ([3]float64, error) {
	raw, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return [3]float64{}, fmt.Errorf("unable to read loads: %w", err)
	}

	var (
		ret   [3]float64
		parts = strings.Split(string(raw), " ")
	)

	for i := range ret {
		ret[i], err = strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return ret, fmt.Errorf("unable to parse item #%d of loadavg: %w", i, err)
		}
	}

	return ret, nil
}
