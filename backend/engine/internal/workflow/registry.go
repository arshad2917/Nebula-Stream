package workflow

import (
	"fmt"
	"sort"
	"sync"
)

type Registry struct {
	mu        sync.RWMutex
	workflows map[string]Definition
	active    string
}

func NewRegistry(initial Definition) *Registry {
	r := &Registry{workflows: make(map[string]Definition)}
	if initial.Name != "" {
		r.workflows[initial.Name] = initial
		r.active = initial.Name
	}
	return r
}

func (r *Registry) Upsert(def Definition) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workflows[def.Name] = def
	if r.active == "" {
		r.active = def.Name
	}
}

func (r *Registry) SetActive(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.workflows[name]; !ok {
		return fmt.Errorf("workflow not found: %s", name)
	}
	r.active = name
	return nil
}

func (r *Registry) Get(name string) (Definition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	def, ok := r.workflows[name]
	return def, ok
}

func (r *Registry) Active() (Definition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	def, ok := r.workflows[r.active]
	return def, ok
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]string, 0, len(r.workflows))
	for name := range r.workflows {
		items = append(items, name)
	}
	sort.Strings(items)
	return items
}
