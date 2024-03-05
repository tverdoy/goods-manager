package cache

import (
	"context"
	"errors"
)

var ErrorNotExists = errors.New("item not found")

// Cache represents an interface for caching operations.
// Implementations are expected to provide methods for getting, setting, and removing items from the cache.
//
//go:generate mockery --name Cache
type Cache interface {
	// Get retrieves a value from the cache based on the provided key.
	// It returns an error if the operation fails.
	Get(ctx context.Context, key string, value interface{}) error

	// Set sets a value in the cache with the provided key.
	// It returns an error if the operation fails.
	Set(ctx context.Context, key string, value interface{}) error

	// Remove removes a value from the cache based on the provided key.
	// It returns an error if the operation fails.
	Remove(ctx context.Context, key string) error
}
