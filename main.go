package main

import (
	"os"

	"github.com/DandaAkhilReddy/claude_peepee/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
