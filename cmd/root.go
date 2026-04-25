package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagConfigPath  string
	flagNoColor     bool
	flagShowAll     bool
	flagInteractive bool
	flagRandom      bool
)

var rootCmd = &cobra.Command{
	Use:   "omcs",
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
	rootCmd.PersistentFlags().StringVar(&flagConfigPath, "config", "", "config file path (default: ~/.config/omcs/config.json)")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disable colored output")
	rootCmd.Flags().BoolVar(&flagShowAll, "all", false, "show all entries including remembered")
	rootCmd.Flags().BoolVarP(&flagInteractive, "interactive", "i", false, "interactive TUI mode")
	rootCmd.Flags().BoolVar(&flagRandom, "random", false, "force a new random shuffle")
}
