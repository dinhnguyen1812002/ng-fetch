package cmd

import (
	"fmt"
	"os"

	"ng-fetch/ascii"
	"ng-fetch/system"

	"github.com/spf13/cobra"
)

var (
	noAscii  bool
	noColors bool
)

var rootCmd = &cobra.Command{
	Use:   "neofetch-go",
	Short: "A simple Neofetch clone written in Go",
	Run: func(cmd *cobra.Command, args []string) {
		runNeofetch()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noAscii, "no-ascii", false, "Disable ASCII art display")
	rootCmd.PersistentFlags().BoolVar(&noColors, "no-colors", false, "Disable colored output")
}

func runNeofetch() {
	// Fetch ASCII art

	if !noAscii {
		ascii.PrintASCIIArt("default") // Fetch ASCII art as a string
	}

	// Print the system info along with ASCII art
	err := system.PrintSystemInfo(noColors)
	if err != nil {
		return
	}
}
