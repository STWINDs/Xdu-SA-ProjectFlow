package cache

import (
	"context"
	"encoding/json"
	"time"

	"cowork/internal/db"
)

func Get(key string, dest interface{}) error {
	if db.Redis == nil {
		return redisUnavailable
	}
	data, err := db.Redis.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func Set(key string, value interface{}, ttl time.Duration) error {
	if db.Redis == nil {
		return redisUnavailable
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.Redis.Set(context.Background(), key, data, ttl).Err()
}

func Del(keys ...string) error {
	if db.Redis == nil {
		return redisUnavailable
	}
	return db.Redis.Del(context.Background(), keys...).Err()
}

func DeletePattern(pattern string) error {
	if db.Redis == nil {
		return redisUnavailable
	}
	iter := db.Redis.Scan(context.Background(), 0, pattern, 0).Iterator()
	for iter.Next(context.Background()) {
		if err := db.Redis.Del(context.Background(), iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

var redisUnavailable = &RedisError{"redis unavailable"}

type RedisError struct{ msg string }

func (e *RedisError) Error() string { return e.msg }
