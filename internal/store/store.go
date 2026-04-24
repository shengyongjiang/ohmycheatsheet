package store

import "github.com/shengyongjiang/ocheetsheet/internal/model"

type StateStore interface {
	GetEntryState(pageKey string, index int) (model.EntryState, bool)
	GetPageStates(pageKey string) map[int]model.EntryState
	SetEntryState(pageKey string, index int, state model.EntryState) error
	ListTrackedPages() []string
	GetDueEntries() []model.DueEntry
	ResetPage(pageKey string) error
	ResetAll() error
	Save() error
}
