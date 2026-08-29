package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rakshitg600/notakto-solo/config"
	"github.com/rakshitg600/notakto-solo/db/generated"
	"github.com/rakshitg600/notakto-solo/store"
)

func loadCoinPackages(ctx context.Context, q *db.Queries) ([]config.CoinPackage, error) {
	value, err := store.GetConfigValueByKey(ctx, q, config.CoinPackagesKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return config.DefaultCoinPackages(), nil
		}
		return nil, fmt.Errorf("get %q config: %w", config.CoinPackagesKey, err)
	}
	packages := make([]config.CoinPackage, 0)
	if err := json.Unmarshal(value, &packages); err != nil {
		return nil, fmt.Errorf("decode %s config: %w", config.CoinPackagesKey, err)
	}
	return packages, nil
}

func loadSignUpConfig(ctx context.Context, q *db.Queries) (config.SignUpConfig, error) {
	value, err := store.GetConfigValueByKey(ctx, q, config.SignUpKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return config.DefaultSignUpConfig(), nil
		}
		return config.SignUpConfig{}, fmt.Errorf("get %q config: %w", config.SignUpKey, err)
	}

	signUp := config.SignUpConfig{}
	if err := json.Unmarshal(value, &signUp); err != nil {
		return config.SignUpConfig{}, fmt.Errorf("decode %s config: %w", config.SignUpKey, err)
	}
	return signUp, nil
}

func loadCensorUsernameConfig(ctx context.Context, q *db.Queries) ([]string, error) {
	value, err := store.GetConfigValueByKey(ctx, q, config.CensorUsernameKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return config.DefaultCensorWords(), nil
		}
		return nil, fmt.Errorf("get %q config: %w", config.CensorUsernameKey, err)
	}

	badWords := make([]string, 0)
	if err := json.Unmarshal(value, &badWords); err != nil {
		return nil, fmt.Errorf("decode %s config: %w", config.CensorUsernameKey, err)
	}
	return badWords, nil
}

func loadUsernameWordLists(ctx context.Context, q *db.Queries) (usernameWordLists, error) {
	adjectives, err := loadUsernameWordList(ctx, q, config.UsernameAdjectivesKey, config.DefaultUsernameAdjectives())
	if err != nil {
		return usernameWordLists{}, err
	}
	animals, err := loadUsernameWordList(ctx, q, config.UsernameAnimalsKey, config.DefaultUsernameAnimals())
	if err != nil {
		return usernameWordLists{}, err
	}
	spaceWords, err := loadUsernameWordList(ctx, q, config.UsernameSpaceWordKey, config.DefaultUsernameSpaceWords())
	if err != nil {
		return usernameWordLists{}, err
	}
	return usernameWordLists{
		adjectives: adjectives,
		animals:    animals,
		spaceWords: spaceWords,
	}, nil
}

func loadUsernameWordList(ctx context.Context, q *db.Queries, key string, fallback []string) ([]string, error) {
	value, err := store.GetConfigValueByKey(ctx, q, key)
	if errors.Is(err, pgx.ErrNoRows) {
		return fallback, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get %q config: %w", key, err)
	}
	words := make([]string, 0)
	if err := json.Unmarshal(value, &words); err != nil {
		return nil, fmt.Errorf("decode %s config: %w", key, err)
	}
	return words, nil
}
