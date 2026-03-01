package utils

import "strings"

// EscapeMarkdown экранирует специальные символы Markdown
func EscapeMarkdown(text string) string {
	// Экранируем специальные символы Markdown
	escaped := text
	escaped = strings.ReplaceAll(escaped, "_", "\\_")
	escaped = strings.ReplaceAll(escaped, "*", "\\*")
	escaped = strings.ReplaceAll(escaped, "[", "\\[")
	escaped = strings.ReplaceAll(escaped, "]", "\\]")
	escaped = strings.ReplaceAll(escaped, "(", "\\(")
	escaped = strings.ReplaceAll(escaped, ")", "\\)")
	escaped = strings.ReplaceAll(escaped, "~", "\\~")
	escaped = strings.ReplaceAll(escaped, "`", "\\`")
	escaped = strings.ReplaceAll(escaped, ">", "\\>")
	escaped = strings.ReplaceAll(escaped, "#", "\\#")
	escaped = strings.ReplaceAll(escaped, "+", "\\+")
	escaped = strings.ReplaceAll(escaped, "-", "\\-")
	escaped = strings.ReplaceAll(escaped, "=", "\\=")
	escaped = strings.ReplaceAll(escaped, "|", "\\|")
	escaped = strings.ReplaceAll(escaped, "{", "\\{")
	escaped = strings.ReplaceAll(escaped, "}", "\\}")
	escaped = strings.ReplaceAll(escaped, ".", "\\.")
	escaped = strings.ReplaceAll(escaped, "!", "\\!")
	return escaped
}

// SplitMessage разбивает длинное сообщение на части
func SplitMessage(text string, maxLength int) []string {
	var parts []string
	currentPart := ""

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(currentPart)+len(line)+1 > maxLength {
			if currentPart != "" {
				parts = append(parts, currentPart)
				currentPart = ""
			}

			// Если одна строка слишком длинная, разбиваем её
			if len(line) > maxLength {
				for len(line) > maxLength {
					parts = append(parts, line[:maxLength])
					line = line[maxLength:]
				}
				if line != "" {
					currentPart = line
				}
			} else {
				currentPart = line
			}
		} else {
			if currentPart != "" {
				currentPart += "\n"
			}
			currentPart += line
		}
	}

	if currentPart != "" {
		parts = append(parts, currentPart)
	}

	return parts
}
