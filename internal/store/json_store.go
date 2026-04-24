package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shengyongjiang/ocheetsheet/internal/model"
)

type stateFile struct {
	SchemaVersion int                        `json:"schema_version"`
	LastModified  time.Time                  `json:"last_modified"`
	Entries       map[string]model.EntryState `json:"entries"`
}

type JSONStore struct {
	path string
	data stateFile
}

func NewJSONStore(path string) (*JSONStore, error) {
	s := &JSONStore{
		path: path,
		data: stateFile{
			SchemaVersion: 1,
			Entries:       make(map[string]model.EntryState),
		},
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, fmt.Errorf("read state file: %w", err)
	}

	if err := json.Unmarshal(raw, &s.data); err != nil {
		return nil, fmt.Errorf("parse state file: %w", err)
	}
	if s.data.Entries == nil {
		s.data.Entries = make(map[string]model.EntryState)
	}
	return s, nil
}

func entryKey(pageKey string, index int) string {
	return fmt.Sprintf("%s/%d", pageKey, index)
}

func parseEntryKey(key string) (string, int, bool) {
	idx := strings.LastIndex(key, "/")
	if idx < 0 {
		return "", 0, false
	}
	n, err := strconv.Atoi(key[idx+1:])
	if err != nil {
		return "", 0, false
	}
	return key[:idx], n, true
}

func (s *JSONStore) GetEntryState(pageKey string, index int) (model.EntryState, bool) {
	es, ok := s.data.Entries[entryKey(pageKey, index)]
	return es, ok
}

func (s *JSONStore) GetPageStates(pageKey string) map[int]model.EntryState {
	result := make(map[int]model.EntryState)
	prefix := pageKey + "/"
	for key, es := range s.data.Entries {
		if strings.HasPrefix(key, prefix) {
			_, idx, ok := parseEntryKey(key)
			if ok {
				result[idx] = es
			}
		}
	}
	return result
}

func (s *JSONStore) SetEntryState(pageKey string, index int, state model.EntryState) error {
	s.data.Entries[entryKey(pageKey, index)] = state
	return nil
}

func (s *JSONStore) ListTrackedPages() []string {
	seen := make(map[string]bool)
	for key := range s.data.Entries {
		page, _, ok := parseEntryKey(key)
		if ok {
			seen[page] = true
		}
	}
	pages := make([]string, 0, len(seen))
	for p := range seen {
		pages = append(pages, p)
	}
	sort.Strings(pages)
	return pages
}

func (s *JSONStore) GetDueEntries() []model.DueEntry {
	var due []model.DueEntry
	now := time.Now()
	for key, es := range s.data.Entries {
		page, idx, ok := parseEntryKey(key)
		if !ok {
			continue
		}
		if es.State == model.StateNeedsReview {
			due = append(due, model.DueEntry{
				PageKey: page,
				Index:   idx,
				State:   es,
			})
			continue
		}
		if es.NextReview != nil && !es.NextReview.After(now) {
			due = append(due, model.DueEntry{
				PageKey: page,
				Index:   idx,
				State:   es,
			})
		}
	}
	return due
}

func (s *JSONStore) ResetPage(pageKey string) error {
	prefix := pageKey + "/"
	for key := range s.data.Entries {
		if strings.HasPrefix(key, prefix) {
			delete(s.data.Entries, key)
		}
	}
	return nil
}

func (s *JSONStore) ResetAll() error {
	s.data.Entries = make(map[string]model.EntryState)
	return nil
}

func (s *JSONStore) Save() error {
	s.data.LastModified = time.Now()
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}
	if err := os.WriteFile(s.path, raw, 0o644); err != nil {
		return fmt.Errorf("write state file: %w", err)
	}
	return nil
}
