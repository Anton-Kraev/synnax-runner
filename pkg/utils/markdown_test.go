package utils

import (
	"strings"
	"testing"
)

func TestEscapeMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escape asterisks",
			input:    "Hello *world*",
			expected: "Hello \\*world\\*",
		},
		{
			name:     "escape underscores",
			input:    "Hello _world_",
			expected: "Hello \\_world\\_",
		},
		{
			name:     "escape brackets",
			input:    "Hello [world](url)",
			expected: "Hello \\[world\\]\\(url\\)",
		},
		{
			name:     "escape backticks",
			input:    "Hello `world`",
			expected: "Hello \\`world\\`",
		},
		{
			name:     "escape multiple symbols",
			input:    "Hello *world* with _emphasis_ and `code`",
			expected: "Hello \\*world\\* with \\_emphasis\\_ and \\`code\\`",
		},
		{
			name:     "no special characters",
			input:    "Hello world",
			expected: "Hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeMarkdown(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitMessage(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxLength int
		expected  []string
	}{
		{
			name:      "short message",
			text:      "Hello world",
			maxLength: 20,
			expected:  []string{"Hello world"},
		},
		{
			name:      "long message",
			text:      "This is a very long message that should be split into multiple parts",
			maxLength: 20,
			expected: []string{
				"This is a very long ",
				"message that should ",
				"be split into multip",
				"le parts",
			},
		},
		{
			name:      "message with newlines",
			text:      "Line 1\nLine 2\nLine 3",
			maxLength: 10,
			expected: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
		},
		{
			name:      "empty message",
			text:      "",
			maxLength: 10,
			expected:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitMessage(tt.text, tt.maxLength)
			if len(result) != len(tt.expected) {
				t.Errorf("SplitMessage() returned %d parts, want %d", len(result), len(tt.expected))
				return
			}
			for i, part := range result {
				if part != tt.expected[i] {
					t.Errorf("SplitMessage() part %d = %q, want %q", i, part, tt.expected[i])
				}
			}
		})
	}
}

func TestSplitMessage_Validation(t *testing.T) {
	// Тест на то, что каждая часть не превышает максимальную длину
	text := strings.Repeat("a", 100)
	maxLength := 20
	
	parts := SplitMessage(text, maxLength)
	
	for i, part := range parts {
		if len(part) > maxLength {
			t.Errorf("Part %d exceeds max length: %d > %d", i, len(part), maxLength)
		}
	}
}
