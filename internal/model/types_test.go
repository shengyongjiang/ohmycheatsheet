package model

import "testing"

func TestMemoryStateNext(t *testing.T) {
	tests := []struct {
		input    MemoryState
		expected MemoryState
	}{
		{StateNotRemembered, StateRemembered},
		{StateRemembered, StateNeedsReview},
		{StateNeedsReview, StateNotRemembered},
	}
	for _, tt := range tests {
		got := tt.input.Next()
		if got != tt.expected {
			t.Errorf("(%s).Next() = %s, want %s", tt.input, got, tt.expected)
		}
	}
}

func TestMemoryStatePrev(t *testing.T) {
	tests := []struct {
		input    MemoryState
		expected MemoryState
	}{
		{StateNotRemembered, StateNeedsReview},
		{StateNeedsReview, StateRemembered},
		{StateRemembered, StateNotRemembered},
	}
	for _, tt := range tests {
		got := tt.input.Prev()
		if got != tt.expected {
			t.Errorf("(%s).Prev() = %s, want %s", tt.input, got, tt.expected)
		}
	}
}

func TestNextPrevRoundTrip(t *testing.T) {
	state := StateNotRemembered
	for range 3 {
		state = state.Next()
	}
	if state != StateNotRemembered {
		t.Errorf("3x Next() should cycle back, got %s", state)
	}

	state = StateNotRemembered
	for range 3 {
		state = state.Prev()
	}
	if state != StateNotRemembered {
		t.Errorf("3x Prev() should cycle back, got %s", state)
	}
}
