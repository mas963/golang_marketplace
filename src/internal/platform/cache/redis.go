package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Cache struct {
	client *redis.Client
}

func Connect(cfg RedisConfig) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Cache{
		client: rdb,
	}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(key, jsonValue, expiration).Err()
}

func (c *Cache) Get(key string, dest interface{}) error {
	val, err := c.client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(key).Err()
}

func (c *Cache) Exists(key string) bool {
	result := c.client.Exists(key)
	return result.Val() > 0
}
