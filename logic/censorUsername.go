package logic

import (
	"errors"
	"strings"
)

var ErrInappropriateUsername = errors.New("username contains inappropriate words")

// ValidateCensorship checks if the given username contains any inappropriate words from the list.
// The matching is case-insensitive and checks if username contains the bad word anywhere (pattern *word*).
func ValidateCensorship(username string, badWords []string) error {
	if len(badWords) == 0 {
		return nil
	}

	lowerUsername := strings.ToLower(username)
	for _, word := range badWords {
		cleanWord := strings.TrimSpace(strings.ToLower(word))
		if cleanWord == "" {
			continue
		}
		if strings.Contains(lowerUsername, cleanWord) {
			return ErrInappropriateUsername
		}
	}
	return nil
}
