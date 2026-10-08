package main

import "testing"

func TestDetectChangeTable(t *testing.T) {
	tests := []struct {
		name, before, after string
		changed             bool
	}{
		{"whitespace is ignored", "price 10", "price\n10", false},
		{"price update alerts", "price 10", "price 11", true},
		{"checkout stock update alerts", "stock available", "stock sold out", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectChange("https://shop.test/item", tt.before, tt.after).Changed
			if got != tt.changed {
				t.Fatalf("changed = %v, want %v", got, tt.changed)
			}
		})
	}
}
