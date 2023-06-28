package embeddingexporter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type persistence interface {
	Persist(logEntry logEntryWithEmbedding) error
}

type RedisPersistence struct {
	rdb *redis.Client
}

func NewRedisPersistence(host string, port string, password string, database int) *RedisPersistence {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       database,
	})
	return &RedisPersistence{rdb: rdb}
}

func (r *RedisPersistence) Persist(logEntry logEntryWithEmbedding) error {
	ctx := context.Background()

	err := r.rdb.HSet(ctx, "myhash", "field1", "value1").Err()
	if err != nil {
		return err
	}

	return nil
}
