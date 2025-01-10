package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"usersvc/pkg/config"
)

type CacheManager interface {
	Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Destroy(ctx context.Context, key string) error
}

type Redis struct {
	Client *redis.Client
}

func NewRedis(config *config.Redis) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Address, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		panic(err)
	}

	log.Println("Redis connection successfully")

	return &Redis{client}
}

func (r *Redis) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	if _, ok := val.(json.Marshaler); ok {
		return r.Client.Set(ctx, key, val, ttl*time.Second).Err()
	}

	val, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, val, ttl*time.Second).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()

	if val != "" && err == nil {
		return val, nil
	}

	if err == redis.Nil {
		return val, errors.New(http.StatusText(http.StatusNotFound))
	}

	return val, err
}

func (r *Redis) Destroy(ctx context.Context, key string) error {
	_, err := r.Client.Do(ctx, "del", key).Result()
	return err
}
