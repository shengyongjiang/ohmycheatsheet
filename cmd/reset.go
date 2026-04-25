package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shengyongjiang/ohmycheatsheet/internal/config"
	"github.com/shengyongjiang/ohmycheatsheet/internal/store"
	"github.com/spf13/cobra"
)

var resetAllFlag bool

var resetCmd = &cobra.Command{
	Use:   "reset [command]",
	Short: "Reset memory state for a command",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runReset,
}

func init() {
	resetCmd.Flags().BoolVar(&resetAllFlag, "all", false, "reset all state")
	rootCmd.AddCommand(resetCmd)
}

func runReset(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && !resetAllFlag {
		return fmt.Errorf("specify a command to reset, or use --all to reset everything")
	}

	cfgPath := flagConfigPath
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	st, err := store.NewJSONStore(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if resetAllFlag {
		fmt.Print("Reset ALL memory state? This cannot be undone. [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(answer)) != "y" {
			fmt.Println("Cancelled.")
			return nil
		}
		if err := st.ResetAll(); err != nil {
			return fmt.Errorf("reset all: %w", err)
		}
		if err := st.Save(); err != nil {
			return fmt.Errorf("save state: %w", err)
		}
		fmt.Println("All state has been reset.")
		return nil
	}

	command := args[0]
	states := st.GetPageStates(command)
	if len(states) == 0 {
		fmt.Printf("No tracked state for %q.\n", command)
		return nil
	}

	fmt.Printf("Reset state for %q (%d entries)? [y/N] ", command, len(states))
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(answer)) != "y" {
		fmt.Println("Cancelled.")
		return nil
	}

	if err := st.ResetPage(command); err != nil {
		return fmt.Errorf("reset page: %w", err)
	}
	if err := st.Save(); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	fmt.Printf("State for %q has been reset.\n", command)
	return nil
}
