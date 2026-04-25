package cmd

import (
	"fmt"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/shengyongjiang/ohmycheatsheet/internal/config"
	"github.com/shengyongjiang/ohmycheatsheet/internal/render"
	"github.com/shengyongjiang/ohmycheatsheet/internal/resolver"
	"github.com/shengyongjiang/ohmycheatsheet/internal/shuffle"
	"github.com/shengyongjiang/ohmycheatsheet/internal/source"
	"github.com/shengyongjiang/ohmycheatsheet/internal/store"
	"github.com/shengyongjiang/ohmycheatsheet/internal/tui"
	"github.com/spf13/cobra"
)

func runShow(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	command := args[0]

	cfgPath := flagConfigPath
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if flagNoColor {
		cfg.ColorEnabled = false
	}

	src := source.NewCheatshSource(cfg.CacheDir)
	res := resolver.New(src)
	page, err := res.Resolve(command)
	if err != nil {
		return fmt.Errorf("command %q not found: %w", command, err)
	}

	st, err := store.NewJSONStore(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	states := st.GetPageStates(page.Name)

	var seed int64
	if flagRandom {
		seed = shuffle.RandomSeed()
		shuffle.SaveSeed(cfg.CacheDir, command, seed)
	} else if saved, err := shuffle.LoadSeed(cfg.CacheDir, command); err == nil {
		seed = saved
	} else {
		seed = shuffle.DailySeed(command)
		shuffle.SaveSeed(cfg.CacheDir, command, seed)
	}

	if flagInteractive {
		m := tui.New(page, states, st, res, seed)
		p := bubbletea.NewProgram(m, bubbletea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	}

	output := render.RenderText(page, states, flagShowAll, !cfg.ColorEnabled, seed)
	fmt.Print(output)
	return nil
}
