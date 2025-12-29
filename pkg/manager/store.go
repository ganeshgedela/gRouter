package manager

import (
	"strings"
	"sync"
)

// ServiceStore manages the registration and retrieval of services.
type ServiceStore struct {
	mu         sync.RWMutex
	serviceMap map[string]Service
}

// NewServiceStore creates a new ServiceStore.
func NewServiceStore() *ServiceStore {
	return &ServiceStore{
		serviceMap: make(map[string]Service),
	}
}

func normalizeService(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Add registers a service with the given name.
func (s *ServiceStore) Add(name string, svc Service) {
	if svc == nil {
		return
	}
	key := normalizeService(name)
	if key == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serviceMap[key] = svc
}

// Get retrieves a service by name.
func (s *ServiceStore) Get(name string) (Service, bool) {
	key := normalizeService(name)
	if key == "" {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	svc, ok := s.serviceMap[key]
	return svc, ok
}

// Delete removes a service by name.
func (s *ServiceStore) Delete(name string) bool {
	key := normalizeService(name)
	if key == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.serviceMap[key]; !ok {
		return false
	}
	delete(s.serviceMap, key)
	return true
}

// Exists checks if a service with the given name is registered.
func (s *ServiceStore) Exists(name string) bool {
	_, ok := s.Get(name)
	return ok
}

// List returns a list of all registered service names.
func (s *ServiceStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.serviceMap))
	for k := range s.serviceMap {
		out = append(out, k)
	}
	return out
}
