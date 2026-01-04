package cmd

import (
	"fmt"
	"runtime"

	"github.com/DandaAkhilReddy/claude_peepee/internal/telemetry"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of claude_peepee",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("claude_peepee version %s\n", telemetry.Version)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
