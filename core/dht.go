package core

import (
	"fmt"
	"sort"
	"sync"
	"crypto/sha1"
	
)

type Node struct {
	ID   string
	Data map[string]string
	mu   sync.RWMutex
}

func NewNode(id string) *Node {
	return &Node{
		ID:   id,
		Data: make(map[string]string),
	}
}

type DHT struct {
	Nodes []*Node
	mu    sync.RWMutex
}

func NewDHT() *DHT {
	return &DHT{
		Nodes: []*Node{},
	}
}

func (d *DHT) AddNode(node *Node) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Nodes = append(d.Nodes, node)
	sort.Slice(d.Nodes, func(i, j int) bool {
		return d.Nodes[i].ID < d.Nodes[j].ID
	})
}

func hashKey(key string) string {
	hash := sha1.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

func (d *DHT) getNodeForKey(key string) *Node {
	if len(d.Nodes) == 0 {
		return nil
	}

	hashedKey := hashKey(key)
	for _, node := range d.Nodes {
		if hashedKey <= node.ID {
			return node
		}
	}

	// If no node is found (wrap around), assign to the first node
	return d.Nodes[0]
}

func (d *DHT) Put(key, value string) {
	node := d.getNodeForKey(key)
	if node == nil {
		fmt.Println("No nodes available to store the key:", key)
		return
	}

	node.mu.Lock()
	defer node.mu.Unlock()
	node.Data[key] = value
	fmt.Printf("Key '%s' stored in Node '%s'\n", key, node.ID)
}

func (d *DHT) Get(key string) (string, bool) {
	node := d.getNodeForKey(key)
	if node == nil {
		fmt.Println("No nodes available to retrieve the key:", key)
		return "", false
	}

	node.mu.RLock()
	defer node.mu.RUnlock()
	value, exists := node.Data[key]
	return value, exists
}

func (d *DHT) Delete(key string) {
	node := d.getNodeForKey(key)
	if node == nil {
		fmt.Println("No nodes available to delete the key:", key)
		return
	}

	node.mu.Lock()
	defer node.mu.Unlock()
	if _, exists := node.Data[key]; exists {
		delete(node.Data, key)
		fmt.Printf("Key '%s' deleted from Node '%s'\n", key, node.ID)
	} else {
		fmt.Printf("Key '%s' not found in Node '%s'\n", key, node.ID)
	}
}

func (d *DHT) PrintDHT() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	fmt.Println("--- DHT State ---")
	for _, node := range d.Nodes {
		node.mu.RLock()
		fmt.Printf("Node ID: %s\n", node.ID)
		fmt.Println("Data:")
		for key, value := range node.Data {
			fmt.Printf("  %s: %s\n", key, value)
		}
		node.mu.RUnlock()
		fmt.Println("------------------")
	}
}

func (d *DHT) GetAllDHTElements() (map[string]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	allElements := make(map[string]string)
	for _, node := range d.Nodes {
		node.mu.RLock()
		for key, value := range node.Data {
			allElements[key] = value
		}
		node.mu.RUnlock()
	}
	return allElements, nil
}