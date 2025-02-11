package core

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)


var dht = NewDHT()

// CombinedPut adds a key-value pair to both the DHT and Redis
func CombinedPut(key, value string) error {
	// Add to DHT
	dht.Put(key, value)

	// Add to Redis
	if err := RedisPut(key, value); err != nil {
		return fmt.Errorf("failed to put in Redis: %v", err)
	}

	dht.PrintDHT()

	return nil
}

// CombinedGet retrieves a value by key, first checking the DHT, then Redis
func CombinedGet(key string) (string, error) {
	// Check DHT
	if value, exists := dht.Get(key); exists {
		return value, nil
	}

	// Check Redis
	value, err := RedisGet(key)
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get from Redis: %v", err)
	}

	// Cache the value in DHT
	dht.Put(key, value)
	
	dht.PrintDHT()

	return value, nil
}

// CombinedDelete removes a key-value pair from both the DHT and Redis
func CombinedDelete(key string) error {
	// Delete from DHT
	dht.Delete(key)

	// Delete from Redis
	if err := RedisDelete(key); err != nil {
		return fmt.Errorf("failed to delete from Redis: %v", err)
	}

	dht.PrintDHT()

	return nil
}

func GetAllElements() (map[string]string, map[string]string, error) {
    // Retrieve DHT elements
    rawDHTElements, err := dht.GetAllDHTElements() // Returns map[string]map[string]string
    if err != nil {
        return nil, nil, fmt.Errorf("failed to retrieve DHT elements: %w", err)
    }

    // Flatten the nested map
    dhtElements := make(map[string]string)
    for nodeID, nodeData := range rawDHTElements {
        for key, value := range nodeData {
            dhtElements[fmt.Sprintf("%s:%s", nodeID, key)] = value // Prefix with node ID
        }
    }

    // Retrieve Redis elements
    redisElements, err := RedisElementsAll() // Replace with your Redis retrieval logic
    if err != nil {
        return nil, nil, fmt.Errorf("failed to retrieve Redis elements: %w", err)
    }

    return dhtElements, redisElements, nil
}
