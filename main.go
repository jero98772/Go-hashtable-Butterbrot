package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)



// DHT represents a simple in-memory distributed hash table
type DHT struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewDHT creates a new DHT instance
func NewDHT() *DHT {
	return &DHT{
		data: make(map[string]string),
	}
}

// Put adds a key-value pair to the DHT
func (d *DHT) Put(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = value
}

// Get retrieves a value by key from the DHT
func (d *DHT) Get(key string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, exists := d.data[key]
	return value, exists
}

// Delete removes a key-value pair from the DHT
func (d *DHT) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.data, key)
}

// PrintDHT prints the current state of the DHT
func (d *DHT) PrintDHT() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	fmt.Println("Current DHT State:")
	for key, value := range d.data {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println()
}

var ctx = context.Background()
var rdb *redis.Client
var dht = NewDHT()

func init() {
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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/api/put", putHandler).Methods("POST")
	r.HandleFunc("/api/get/{key}", getHandler).Methods("GET")
	r.HandleFunc("/api/delete/{key}", deleteHandler).Methods("DELETE")

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for key, value := range data {
		if err := CombinedPut(key, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data inserted successfully"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := CombinedGet(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{key: value})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := CombinedDelete(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key deleted successfully"))
}