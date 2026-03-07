package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type VacancyCache struct {
	redisClient *redis.Client
}

func NewVacancyCache(redisClient *redis.Client) *VacancyCache {
	return &VacancyCache{redisClient: redisClient}
}

func VacancyKey(hhID string) string {
	return "vacancy:" + hhID
}
func (v *VacancyCache) SetSeen(ctx context.Context, hhID string) error {
	_, err := v.redisClient.SetNX(ctx, VacancyKey(hhID), 1, 24*time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}

func (v *VacancyCache) IsSeen(ctx context.Context, hhID string) (bool, error) {
	count, err := v.redisClient.Exists(ctx, VacancyKey(hhID)).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
