package main

import (
	"os"

	"github.com/DandaAkhilReddy/claude_pp/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
