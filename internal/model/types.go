package model

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type MemoryState string

const (
	StateNotRemembered MemoryState = "not_remembered"
	StateRemembered    MemoryState = "remembered"
	StateNeedsReview   MemoryState = "needs_review"
)

func (s MemoryState) Next() MemoryState {
	switch s {
	case StateNotRemembered:
		return StateRemembered
	case StateRemembered:
		return StateNeedsReview
	case StateNeedsReview:
		return StateNotRemembered
	default:
		return StateNotRemembered
	}
}

func (s MemoryState) Prev() MemoryState {
	switch s {
	case StateNotRemembered:
		return StateNeedsReview
	case StateNeedsReview:
		return StateRemembered
	case StateRemembered:
		return StateNotRemembered
	default:
		return StateNotRemembered
	}
}

type Entry struct {
	Index       int
	Description string
	Command     string
	Fingerprint string
}

type Page struct {
	Name        string
	Platform    string
	Description string
	SeeAlso     []string
	URL         string
	Entries     []Entry
}

type EntryState struct {
	State        MemoryState `json:"state"`
	LastReviewed  time.Time   `json:"last_reviewed,omitempty"`
	ReviewCount   int         `json:"review_count"`
	Fingerprint   string      `json:"fingerprint"`
	EaseFactor   float64     `json:"ease_factor,omitempty"`
	IntervalDays int         `json:"interval_days,omitempty"`
	NextReview   *time.Time  `json:"next_review,omitempty"`
}

type DueEntry struct {
	PageKey string
	Index   int
	Entry   Entry
	State   EntryState
}

func ComputeFingerprint(description, command string) string {
	normalized := strings.ToLower(strings.TrimSpace(description)) + "\n" + strings.ToLower(strings.TrimSpace(command))
	hash := sha256.Sum256([]byte(normalized))
	return fmt.Sprintf("%x", hash)
}
