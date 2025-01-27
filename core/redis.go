package core

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	// Add nodes to DHT
	dht.AddNode(NewNode(hashKey("Node1")))
	dht.AddNode(NewNode(hashKey("Node2")))
	dht.AddNode(NewNode(hashKey("Node3")))

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

// RedisPut adds a key-value pair to Redis
func RedisPut(key, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

// RedisGet retrieves a value by key from Redis
func RedisGet(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// RedisDelete removes a key-value pair from Redis
func RedisDelete(key string) error {
	return rdb.Del(ctx, key).Err()
}

func RedisElementsAll() (map[string]string, error) {
    // Fetch all keys using the Redis client
    keys, err := rdb.Keys(ctx, "*").Result()
    if err != nil {
        return nil, err
    }

    result := make(map[string]string)
    for _, key := range keys {
        value, err := rdb.Get(ctx, key).Result()
        if err != nil {
            continue // Skip keys that can't be fetched
        }
        result[key] = value
    }
    return result, nil
}
