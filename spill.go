package spill

import (
	"sync"
)

// SpillMap represents a collection of maps with limited capacity per map
type SpillMap struct {
	maps       []map[interface{}]interface{}
	maxPerMap  int
	currentMap int
	mutex      sync.RWMutex
}

// NewSpillMap creates a new SpillMap with the specified key limit per map
func NewSpillMap(maxKeysPerMap int) *SpillMap {
	if maxKeysPerMap <= 0 {
		maxKeysPerMap = 100 // Default value
	}

	return &SpillMap{
		maps:       []map[interface{}]interface{}{make(map[interface{}]interface{})},
		maxPerMap:  maxKeysPerMap,
		currentMap: 0,
		mutex:      sync.RWMutex{},
	}
}

// Put adds a key-value pair to the SpillMap
func (sm *SpillMap) Put(key, value interface{}) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if current map has reached capacity
	if len(sm.maps[sm.currentMap]) >= sm.maxPerMap {
		// Create a new map
		sm.maps = append(sm.maps, make(map[interface{}]interface{}))
		sm.currentMap++
	}

	// Add the key-value pair to the current map
	sm.maps[sm.currentMap][key] = value
}

// Get retrieves a value for a key from the SpillMap
func (sm *SpillMap) Get(key interface{}) (interface{}, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Search all maps for the key
	for _, m := range sm.maps {
		if val, found := m[key]; found {
			return val, true
		}
	}

	return nil, false
}

// Delete removes a key-value pair from the SpillMap
func (sm *SpillMap) Delete(key interface{}) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Search all maps for the key
	for i, m := range sm.maps {
		if _, found := m[key]; found {
			delete(m, key)

			// If map becomes empty and it's not the last map, remove it
			if len(m) == 0 && len(sm.maps) > 1 && i != sm.currentMap {
				sm.maps = append(sm.maps[:i], sm.maps[i+1:]...)
				if i <= sm.currentMap {
					sm.currentMap--
				}
			}

			return true
		}
	}

	return false
}

// Size returns the total number of key-value pairs in the SpillMap
func (sm *SpillMap) Size() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	total := 0
	for _, m := range sm.maps {
		total += len(m)
	}
	return total
}

// MapCount returns the number of maps currently in use
func (sm *SpillMap) MapCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return len(sm.maps)
}

// ForEach iterates over all key-value pairs in the SpillMap
func (sm *SpillMap) ForEach(fn func(key, value interface{}) bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for _, m := range sm.maps {
		for k, v := range m {
			if !fn(k, v) {
				return
			}
		}
	}
}
