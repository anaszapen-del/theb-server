package container

import (
	"fmt"
	"sync"
)

// Container is a simple dependency injection container
type Container struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// New creates a new container
func New() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

// Register registers a service in the container
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get retrieves a service from the container
func (c *Container) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// Has checks if a service exists in the container
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	return exists
}

// GetTyped retrieves a service with type assertion
func GetTyped[T any](c *Container, name string) (T, error) {
	var zero T
	service, err := c.Get(name)
	if err != nil {
		return zero, err
	}

	typed, ok := service.(T)
	if !ok {
		return zero, fmt.Errorf("service %s is not of expected type", name)
	}

	return typed, nil
}
