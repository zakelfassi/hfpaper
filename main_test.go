package main

import "testing"

func TestParsePaperID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2602.08025", "2602.08025"},
		{"2602.08025v1", "2602.08025v1"},
		{"https://huggingface.co/papers/2602.08025", "2602.08025"},
		{"https://huggingface.co/papers/2602.08025.md", "2602.08025"},
		{"https://arxiv.org/abs/2602.08025", "2602.08025"},
		{"https://arxiv.org/pdf/2602.08025", "2602.08025"},
		{"2503.12345", "2503.12345"},
		{"https://huggingface.co/papers/2503.12345v2", "2503.12345v2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parsePaperID(tt.input)
			if got != tt.expected {
				t.Errorf("parsePaperID(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
