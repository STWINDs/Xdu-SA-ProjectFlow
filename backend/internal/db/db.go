package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Redis *redis.Client

func SetDB(d *gorm.DB) {
	DB = d
}

func SetRedis(r *redis.Client) {
	Redis = r
}

func GetRedis() *redis.Client {
	return Redis
}

func RedisPing() error {
	return Redis.Ping(context.Background()).Err()
}

func InitRedis(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v (will fallback to DB)", err)
		return nil
	}
	log.Println("Redis connected successfully")
	return client
}
