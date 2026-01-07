package storage

import (
	"strings"
	"unicode"
)

func sanitizePlayerName(name string) string {
	var result strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-' {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) {
			result.WriteRune('_')
		}
	}
	sanitized := result.String()
	if sanitized == "" {
		return "Player"
	}
	return sanitized
}
