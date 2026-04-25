package cmd

import (
	"fmt"

	"github.com/shengyongjiang/ohmycheatsheet/internal/config"
	"github.com/shengyongjiang/ohmycheatsheet/internal/model"
	"github.com/shengyongjiang/ohmycheatsheet/internal/resolver"
	"github.com/shengyongjiang/ohmycheatsheet/internal/source"
	"github.com/shengyongjiang/ohmycheatsheet/internal/store"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [command]",
	Short: "Show memory statistics",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
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

	src := source.NewCheatshSource(cfg.CacheDir)
	res := resolver.New(src)

	var pagesToShow []string
	if len(args) > 0 {
		pagesToShow = []string{args[0]}
	} else {
		pagesToShow = st.ListTrackedPages()
	}

	if len(pagesToShow) == 0 {
		fmt.Println("No tracked commands yet. Use `omcs <command> -i` to start learning.")
		return nil
	}

	totalEntries := 0
	totalRemembered := 0
	totalReview := 0
	totalNotRemembered := 0

	for _, pageKey := range pagesToShow {
		page, err := res.Resolve(pageKey)
		if err != nil {
			continue
		}
		states := st.GetPageStates(pageKey)

		remembered := 0
		review := 0
		notRemembered := 0

		for _, entry := range page.Entries {
			es, ok := states[entry.Index]
			if !ok || es.State == model.StateNotRemembered {
				notRemembered++
			} else if es.State == model.StateRemembered {
				remembered++
			} else if es.State == model.StateNeedsReview {
				review++
			}
		}

		total := len(page.Entries)
		pct := 0
		if total > 0 {
			pct = remembered * 100 / total
		}

		fmt.Printf("  %-20s %d/%d remembered (%d%%)  %d to review  %d remaining\n",
			pageKey, remembered, total, pct, review, notRemembered)

		totalEntries += total
		totalRemembered += remembered
		totalReview += review
		totalNotRemembered += notRemembered
	}

	if len(pagesToShow) > 1 {
		totalPct := 0
		if totalEntries > 0 {
			totalPct = totalRemembered * 100 / totalEntries
		}
		fmt.Println()
		fmt.Printf("  %-20s %d/%d remembered (%d%%)  %d to review  %d remaining\n",
			"TOTAL", totalRemembered, totalEntries, totalPct, totalReview, totalNotRemembered)
	}

	return nil
}
