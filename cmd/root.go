package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagConfigPath  string
	flagTldrPath    string
	flagNoColor     bool
	flagVerbose     bool
	flagShowAll     bool
	flagInteractive bool
)

var rootCmd = &cobra.Command{
	Use:   "ocs",
	Short: "Cheatsheet with memory",
	Long:  "A cheatsheet tool that tracks which commands you've memorized.",
	Args:  cobra.ArbitraryArgs,
	RunE:  runShow,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagConfigPath, "config", "", "config file path (default: ~/.config/ocs/config.json)")
	rootCmd.PersistentFlags().StringVar(&flagTldrPath, "tldr-path", "", "override tldr cache path")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "debug output")
	rootCmd.Flags().BoolVar(&flagShowAll, "all", false, "show all entries including remembered")
	rootCmd.Flags().BoolVarP(&flagInteractive, "interactive", "i", false, "interactive TUI mode")
}
