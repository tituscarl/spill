package spill

import (
	"testing"
)

func TestSpillMap(t *testing.T) {
	// Create a map with max 5 keys per internal map
	sm := NewSpillMap(5)

	// Test initial state
	if sm.MapCount() != 1 {
		t.Errorf("Expected 1 map, got %d", sm.MapCount())
	}

	// Add 10 items, should create 2 maps
	for i := 0; i < 10; i++ {
		sm.Put(i, i*10)
	}

	if sm.MapCount() != 2 {
		t.Errorf("Expected 2 maps, got %d", sm.MapCount())
	}

	// Check that we can retrieve all values
	for i := 0; i < 10; i++ {
		val, found := sm.Get(i)
		if !found {
			t.Errorf("Key %d not found", i)
			continue
		}

		if val.(int) != i*10 {
			t.Errorf("Expected value %d, got %v", i*10, val)
		}
	}

	// Test size
	if sm.Size() != 10 {
		t.Errorf("Expected size 10, got %d", sm.Size())
	}

	// Test deletion
	sm.Delete(1)
	sm.Delete(6)

	if sm.Size() != 8 {
		t.Errorf("Expected size 8 after deletion, got %d", sm.Size())
	}

	_, found := sm.Get(1)
	if found {
		t.Errorf("Key 1 should have been deleted")
	}

	// Test ForEach
	count := 0
	sm.ForEach(func(key, value interface{}) bool {
		count++
		return true
	})

	if count != 8 {
		t.Errorf("ForEach should have iterated over 8 items, counted %d", count)
	}

	// Test early termination of ForEach
	count = 0
	sm.ForEach(func(key, value interface{}) bool {
		count++
		return count < 3 // Stop after processing 3 items
	})

	if count != 3 {
		t.Errorf("ForEach should have stopped after 3 iterations, counted %d", count)
	}
}

func TestSpillMapTypeSpecific(t *testing.T) {
	// Test with string keys and various value types
	sm := NewSpillMap(3)

	sm.Put("name", "John")
	sm.Put("age", 30)
	sm.Put("employed", true)
	sm.Put("salary", 50000.50)

	if sm.MapCount() != 2 {
		t.Errorf("Expected 2 maps, got %d", sm.MapCount())
	}

	if sm.Size() != 4 {
		t.Errorf("Expected size 4, got %d", sm.Size())
	}

	// Verify values
	name, _ := sm.Get("name")
	if name.(string) != "John" {
		t.Errorf("Expected 'John', got %v", name)
	}

	employed, _ := sm.Get("employed")
	if !employed.(bool) {
		t.Errorf("Expected true, got %v", employed)
	}
}
