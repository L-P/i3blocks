package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("usage: %s (bandwidth|loadavg|ups)\n", os.Args[0])
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case "loadavg":
		if err := loadavg(); err != nil {
			panic(fmt.Errorf("unable to compute loadavg: %w", err))
		}
	case "bandwidth":
		if err := bandwidth(); err != nil {
			panic(fmt.Errorf("unable to compute bandwidth: %w", err))
		}
	case "ups":
		if err := ups(); err != nil {
			panic(fmt.Errorf("unable to obtain ups info: %w", err))
		}
	default:
		panic(fmt.Errorf("unknown subcommand: %s", flag.Args()[0]))
	}
}
