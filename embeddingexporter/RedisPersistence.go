package embeddingexporter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisPersistence struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisPersistence(host string, port string, password string, database int) *RedisPersistence {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       database,
	})
	return &RedisPersistence{
		rdb: rdb,
		ttl: time.Hour * 24 * 30,
	}
}

func (r *RedisPersistence) Persist(key string, properties Properties) error {
	ctx := context.Background()

	error := r.rdb.HSet(ctx, key, map[string]interface{}(properties)).Err()
	if error != nil {
		return error
	}
	error = r.rdb.Expire(ctx, key, r.ttl).Err()
	if error != nil {
		return error
	}

	return nil
}
