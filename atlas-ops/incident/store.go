package incident

import (
	"fmt"
	"sync"
)

var (
	store   = make(map[string]*Incident)
	mu      sync.Mutex
	counter int
)

func Create(i *Incident) string {
	mu.Lock()
	defer mu.Unlock()

	counter++
	id := fmt.Sprintf("INC-%d", counter)
	i.ID = id
	store[id] = i
	return id
}

func Get(id string) (*Incident, bool) {
	mu.Lock()
	defer mu.Unlock()

	inc, ok := store[id]
	return inc, ok
}