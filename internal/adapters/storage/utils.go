package storage

import "strings"

func sanitizePlayerName(name string) string {
	if name == "" {
		return "Player"
	}
	return strings.ReplaceAll(name, " ", "_")
}
