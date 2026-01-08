package manager

import (
	"fmt"
	"sync"
)

// Factory is a function that creates a Service instance using the provided dependencies.
type Factory func(deps Deps) (Service, error)

var (
	// registryMu protects the factory registry map.
	registryMu sync.RWMutex
	// factories holds the registered service factories.
	factories = make(map[string]Factory)
)

// RegisterFactory registers a factory function for a given service ID.
// This function is typically called in the init() function of service packages.
// It panics if a factory with the same ID is already registered to fail fast at startup.
func RegisterFactory(id string, f Factory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := factories[id]; exists {
		panic(fmt.Sprintf("service factory already registered for id: %s", id))
	}
	factories[id] = f
}

// GetFactory retrieves a registered factory by service ID.
func GetFactory(id string) (Factory, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	f, ok := factories[id]
	return f, ok
}
