package cache

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	c := NewCache(time.Second)

	// Test adding and getting an item
	c.Add("key1", "value1", 2*time.Second)
	value, found := c.Get("key1")
	if !found || value != "value1" {
		t.Errorf("Expected to find 'value1' for 'key1', got %v", value)
	}

	// Test getting a non-existent item
	value, found = c.Get("key2")
	if found {
		t.Errorf("Expected not to find 'key2'")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache(time.Second)

	// Add an item
	c.Add("key1", "value1", 2*time.Second)

	// Delete the item
	c.Delete("key1")

	// Check if the item is deleted
	value, found := c.Get("key1")
	if found {
		t.Errorf("Expected 'key1' to be deleted, got %v", value)
	}
}

func TestCache_Clear(t *testing.T) {
	c := NewCache(time.Second)

	// Add items
	c.Add("key1", "value1", 2*time.Second)
	c.Add("key2", "value2", 2*time.Second)

	// Clear the cache
	c.Clear()

	// Check if all items are cleared
	value, found := c.Get("key1")
	if found {
		t.Errorf("Expected 'key1' to be cleared, got %v", value)
	}

	value, found = c.Get("key2")
	if found {
		t.Errorf("Expected 'key2' to be cleared, got %v", value)
	}
}

func TestCache_Expiry(t *testing.T) {
	c := NewCache(100 * time.Millisecond)

	// Add an item with a short expiry
	c.Add("key1", "value1", 100*time.Millisecond)

	// Wait for the item to expire
	time.Sleep(200 * time.Millisecond)

	// Check if the item is expired
	value, found := c.Get("key1")
	if found {
		t.Errorf("Expected 'key1' to be expired, got %v", value)
	}
}
