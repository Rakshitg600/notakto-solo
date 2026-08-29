package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	db "github.com/rakshitg600/notakto-solo/db/generated"
	"github.com/rakshitg600/notakto-solo/logic"
	"github.com/rakshitg600/notakto-solo/store"
)

const maxUsernameCollisionSuffix = 9999

type usernameWordLists struct {
	adjectives []string
	animals    []string
	spaceWords []string
}

func generateAvailableUsername(ctx context.Context, q *db.Queries, lists usernameWordLists) (string, error) {
	base, err := generateUsernameBase(lists)
	if err != nil {
		return "", err
	}
	if err := logic.ValidateUsername(base); err != nil {
		return "", fmt.Errorf("generated username is invalid: %w", err)
	}

	exists, err := store.CheckUsernameExists(ctx, q, base)
	if err != nil {
		return "", err
	}
	if !exists {
		return base, nil
	}

	for suffix := 1; suffix <= maxUsernameCollisionSuffix; suffix++ {
		candidate := fmt.Sprintf("%s%d", base, suffix)
		if err := logic.ValidateUsername(candidate); err != nil {
			return "", fmt.Errorf("generated username suffix exceeded validation limits: %w", err)
		}

		exists, err := store.CheckUsernameExists(ctx, q, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not generate available username after %d suffix attempts", maxUsernameCollisionSuffix)
}

func generateUsernameBase(lists usernameWordLists) (string, error) {
	if len(lists.adjectives) == 0 {
		return "", errors.New("username adjectives list cannot be empty")
	}
	if len(lists.animals) == 0 {
		return "", errors.New("username animals list cannot be empty")
	}
	if len(lists.spaceWords) == 0 {
		return "", errors.New("username space words list cannot be empty")
	}

	adjective := lists.adjectives[rand.Int()%len(lists.adjectives)]
	animal := lists.animals[rand.Int()%len(lists.animals)]
	spaceWord := lists.spaceWords[rand.Int()%len(lists.spaceWords)]

	return fmt.Sprintf("%s-%s-%s", adjective, animal, spaceWord), nil
}
