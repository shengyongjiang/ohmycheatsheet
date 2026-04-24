package cmd

import (
	"fmt"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/shengyongjiang/ocheetsheet/internal/config"
	"github.com/shengyongjiang/ocheetsheet/internal/parser"
	"github.com/shengyongjiang/ocheetsheet/internal/render"
	"github.com/shengyongjiang/ocheetsheet/internal/resolver"
	"github.com/shengyongjiang/ocheetsheet/internal/store"
	"github.com/shengyongjiang/ocheetsheet/internal/tui"
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
	if flagTldrPath != "" {
		cfg.TldrCachePath = flagTldrPath
	}
	if flagNoColor {
		cfg.ColorEnabled = false
	}

	res := resolver.NewDefault(cfg.TldrCachePath)
	path, err := res.Resolve(command)
	if err != nil {
		return fmt.Errorf("command %q not found. Make sure tldr cache is populated (run: tldr --update)", command)
	}

	page, err := parser.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	st, err := store.NewJSONStore(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	states := st.GetPageStates(page.Name)

	if flagInteractive {
		m := tui.New(page, states, st, res)
		p := bubbletea.NewProgram(m, bubbletea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	}

	output := render.RenderText(page, states, flagShowAll, !cfg.ColorEnabled)
	fmt.Print(output)
	return nil
}
