package godi

import (
	"fmt"
	"sync"
)

// Lifecycle determines how a service is instantiated and cached
type Lifecycle int

const (
	// Transient creates a new instance each time the service is resolved
	Transient Lifecycle = iota
	// Singleton creates only one instance and reuses it for all subsequent resolves
	Singleton
)

// serviceEntry represents a service in the container
type serviceEntry[T any] struct {
	constructor func() T
	instance    T
	created     bool
	lifecycle   Lifecycle
	mu          sync.Mutex
}

// Container represents the DI container
type Container struct {
	mu       sync.RWMutex
	services map[string]any
}

// New creates a new DI container
func New() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

// Register adds a service to the container
func Register[T any](c *Container, constructor func() T, lifecycle Lifecycle) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%T", *new(T))
	c.services[key] = &serviceEntry[T]{
		constructor: constructor,
		lifecycle:   lifecycle,
	}
}

// Resolve retrieves a service from the container
func Resolve[T any](c *Container) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%T", *new(T))
	entry, exists := c.services[key]
	if !exists {
		var zero T
		return zero, fmt.Errorf("service of type %s not registered", key)
	}

	if serviceEntry, ok := entry.(*serviceEntry[T]); ok {
		if serviceEntry.lifecycle == Singleton {
			serviceEntry.mu.Lock()
			defer serviceEntry.mu.Unlock()
			if !serviceEntry.created {
				serviceEntry.instance = serviceEntry.constructor()
				serviceEntry.created = true
			}
			return serviceEntry.instance, nil
		}
		return serviceEntry.constructor(), nil
	}

	var zero T
	return zero, fmt.Errorf("invalid service entry for type %s", key)
}

// MustResolve retrieves a service from the container or panics if not found
func MustResolve[T any](c *Container) T {
	service, err := Resolve[T](c)
	if err != nil {
		panic(err)
	}
	return service
}
