package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shengyongjiang/ocheetsheet/internal/model"
	"github.com/shengyongjiang/ocheetsheet/internal/parser"
	"github.com/shengyongjiang/ocheetsheet/internal/render"
	"github.com/shengyongjiang/ocheetsheet/internal/resolver"
	"github.com/spf13/cobra"
)

func runShow(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	command := args[0]
	cachePath := flagTldrPath
	if cachePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		cachePath = filepath.Join(home, ".tldr", "cache", "pages")
	}

	res := resolver.NewDefault(cachePath)
	path, err := res.Resolve(command)
	if err != nil {
		return fmt.Errorf("command %q not found. Make sure tldr cache is populated (run: tldr --update)", command)
	}

	page, err := parser.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	states := map[int]model.EntryState{}

	output := render.RenderText(page, states, flagShowAll, flagNoColor)
	fmt.Print(output)
	return nil
}
