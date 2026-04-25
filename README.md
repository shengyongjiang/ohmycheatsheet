# omcs — Oh My Cheat Sheet

A terminal cheatsheet tool that tracks which commands you've memorized. Powered by [cheat.sh](https://cheat.sh) as data source.

## Install

```bash
go install github.com/shengyongjiang/ohmycheatsheet@latest
```

Or build from source:

```bash
git clone https://github.com/shengyongjiang/ohmycheatsheet.git
cd ohmycheatsheet
go build -o omcs .
```

## Usage

### View a cheatsheet

```bash
omcs git          # show 10 random entries you haven't memorized
omcs git --all    # show all entries
omcs git --random # reshuffle entries randomly
```

The default view shows 10 non-remembered entries in a randomized order that stays consistent across runs. Use `--random` to force a new shuffle.

### Interactive mode

```bash
omcs git -i
```

Opens a TUI with all entries. The entry order matches the last non-interactive output, so you can preview with `omcs git` then drill in with `omcs git -i`.

#### Key bindings

| Key | Action |
|---|---|
| `j` / `k` / Arrow keys | Navigate entries |
| `Left` / `Right` | Cycle memory state |
| `x` / `X` | Mark as remembered |
| `Enter` | Mark as needs review |
| `a` | Toggle show all / filter |
| `r` | Reset all states |
| `q` / `Esc` | Quit (auto-saves) |
| `?` | Help |

### Memory states

- **not remembered** (`o`, white) — default, always shown
- **remembered** (`x`, gray) — hidden by default, means you know this one
- **needs review** (`+`, red) — highlighted, you want to revisit this

### Other commands

```bash
omcs review          # flashcard-style review of "needs review" entries
omcs review git      # review only git entries
omcs stats           # show memorization progress across all commands
omcs stats git       # show progress for a specific command
omcs list            # list all tracked commands
omcs reset git       # reset memory state for git
omcs reset --all     # reset everything
omcs completion zsh  # generate shell completions (bash/zsh/fish)
```

## How it works

1. Fetches cheatsheet content from `cheat.sh` on first use, caches locally for 7 days
2. Entries are shuffled with a deterministic daily seed (consistent within a day)
3. Memory state is persisted locally as JSON
4. When entries are marked as remembered and hidden, related entries (e.g. `git-log`, `git-stash` for `git`) backfill the list

## Data storage

| Data | Location (macOS) |
|---|---|
| Cheatsheet cache | `~/Library/Caches/omcs/cheatsh/` |
| Memory state | `~/Library/Application Support/omcs/state.json` |
| Config (optional) | `~/Library/Application Support/omcs/config.json` |
| Shuffle seed | `~/Library/Caches/omcs/cheatsh/seeds/` |

On Linux, paths follow XDG defaults (`~/.cache/omcs/` and `~/.config/omcs/`).

## Project structure

```
cmd/              CLI commands (show, review, stats, list, reset, completion)
internal/
  config/         Configuration loading
  model/          Core types (Entry, Page, MemoryState, EntryState)
  parser/         Parses cheat.sh plain-text format into structured entries
  render/         Non-interactive text renderer
  resolver/       Resolves command names to pages via cheat.sh source
  shuffle/        Deterministic entry shuffling with seed persistence
  source/         HTTP client for cheat.sh with local caching
  store/          JSON-based state persistence
  tui/            Interactive TUI (Bubble Tea)
```

## Language

[简体中文](README_CN.md)
