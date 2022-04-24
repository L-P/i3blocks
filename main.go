package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("usage: %s (bandwidth)\n", os.Args[0])
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case "bandwidth":
		if err := bandwidth(); err != nil {
			panic(fmt.Errorf("unable to compute bandwidth: %w", err))
		}
	default:
		panic(fmt.Errorf("unknown subcommand: %s", flag.Args()[0])) // nolint:goerr113
	}
}
