package blockchain

import "testing"

func TestIsStuckAmountEnough(t *testing.T) {
	tests := []struct {
		expected string
		received string
		want     bool
	}{
		{"100", "5", true},
		{"100", "4.99", false},
		{"100", "100", true},
		{"100", "0", false},
		{"0.01", "0.0005", true},
		{"0.01", "0.000499999999999999", false},
		{"", "0.000001", true},  // expected amount unknown: check skipped
		{"0", "0.000001", true}, // expected amount unknown: check skipped
		{"abc", "1", true},      // invalid expected amount: check skipped
		{"100", "", false},
	}
	for _, tt := range tests {
		a := Address{Amount: tt.expected}
		if got := a.IsStuckAmountEnough(tt.received); got != tt.want {
			t.Errorf("expected=%q received=%q: got %v, want %v", tt.expected, tt.received, got, tt.want)
		}
	}
}
