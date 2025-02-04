package redis

import (
	"context"
	"encoding/json"
	"shopping-gamification/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redis *redis.Client
}

func NewRepository(redis *redis.Client) *Repository {
	return &Repository{redis: redis}
}

func (r *Repository) GetPageHome(ctx context.Context) (domain.PageHome, error) {
	var pageHome domain.PageHome

	val, err := r.redis.Get(ctx, "page_home").Result()
	if err == redis.Nil {
		return pageHome, nil // Cache miss, return empty result
	} else if err != nil {
		return pageHome, err
	}

	err = json.Unmarshal([]byte(val), &pageHome)
	if err != nil {
		return pageHome, err
	}

	return pageHome, nil
}

func (r *Repository) SetPageHome(ctx context.Context, pageHome domain.PageHome) error {
	data, err := json.Marshal(pageHome)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, "page_home", data, 0).Err()
}

func (r *Repository) GetRedisValue(ctx context.Context, key string) (string, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Repository) SetRedisValue(ctx context.Context, key string, value string) error {
	return r.redis.Set(ctx, key, value, 0).Err()
}
