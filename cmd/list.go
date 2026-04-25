package cmd

import (
	"fmt"

	"github.com/shengyongjiang/ohmycheatsheet/internal/config"
	"github.com/shengyongjiang/ohmycheatsheet/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all commands with tracked state",
	Args:  cobra.NoArgs,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
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

	pages := st.ListTrackedPages()
	if len(pages) == 0 {
		fmt.Println("No tracked commands yet. Use `omcs <command> -i` to start learning.")
		return nil
	}

	fmt.Println("Tracked commands:")
	for _, p := range pages {
		fmt.Printf("  %s\n", p)
	}
	return nil
}
