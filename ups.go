package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func ups() error {
	state, err := upsc("eaton3s700")
	if err != nil {
		return fmt.Errorf("unable to obtain UPS state: %w", err)
	}

	secondsLeft, err := strconv.Atoi(state["battery.runtime"])
	if err != nil {
		return fmt.Errorf("unable to parse battery runtime: %w", err)
	}

	left := time.Duration(secondsLeft) * time.Second

	fmt.Printf("%s%% (%s)\n", state["battery.charge"], left)

	return nil
}

func upsc(name string) (map[string]string, error) {
	out, err := exec.Command("upsc", name).Output()
	if err != nil {
		return nil, fmt.Errorf("unable to run upsc: %w", err)
	}

	state := make(map[string]string, 46 /* wc -l of a local run */)
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		state[parts[0]] = parts[1]
	}

	return state, nil
}
