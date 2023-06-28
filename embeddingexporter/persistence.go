package embeddingexporter

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type persistence interface {
	Persist(logEntry logEntryWithEmbedding) error
}

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

func (r *RedisPersistence) Persist(logEntry logEntryWithEmbedding) error {
	ctx := context.Background()

	//Probably want to grab the system from the context and add it to the key
	id := fmt.Sprintf("log_%s_%s_%s",
		logEntry.logEntry.level,
		logEntry.logEntry.TraceId,
		uuid.New().String())
	bytes, error := float32SliceToByteSlice(logEntry.embedding)

	if error != nil {
		return error
	}

	properties := map[string]interface{}{
		"timestamp": logEntry.logEntry.timestamp.Unix(),
		"body":      logEntry.logEntry.body,
		"embedding": bytes,
		"level":     logEntry.logEntry.level,
		"traceId":   logEntry.logEntry.TraceId,
		"spanId":    logEntry.logEntry.SpanId,
	}

	error = r.rdb.HSet(ctx, id, properties).Err()
	if error != nil {
		return error
	}
	error = r.rdb.Expire(ctx, id, r.ttl).Err()
	if error != nil {
		return error
	}

	return nil
}

func float32SliceToByteSlice(floats []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, floats)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
