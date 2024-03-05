package domain

import (
	"context"
	"errors"
	"goods-manager/internal/domain/entity"
)

var ErrorGoodNotFound = errors.New("good not found")

// GoodUsecase represents the use case interface for managing goods.
//
//go:generate mockery --name GoodUsecase
type GoodUsecase interface {
	// Create creates a new Good entity.
	Create(ctx context.Context, good *entity.Good) error

	// Get retrieves a Good entity by its ID.
	Get(ctx context.Context, id int) (*entity.Good, error)

	// Update updates an existing Good entity.
	Update(ctx context.Context, good *entity.Good) error

	// Delete deletes an existing Good entity.
	Delete(ctx context.Context, good *entity.Good) error

	// List retrieves a list of Good entities with pagination support.
	List(ctx context.Context, limit, offset int) ([]*entity.Good, error)

	// Reprioritize changes the priority of a Good entity identified by its ID.
	// It takes a context.Context, ID of the Good, and a new priority as parameters.
	// It returns a map containing IDs of affected Goods and their new priorities, and an error if the operation fails.
	Reprioritize(ctx context.Context, id, newPriority int) (map[int]int, error)
}

//go:generate mockery --name GoodRepository
type GoodRepository interface {
	Create(ctx context.Context, good *entity.Good) error
	Get(ctx context.Context, id int) (*entity.Good, error)
	Update(ctx context.Context, good *entity.Good) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entity.Good, error)

	// Reprioritize changes the priority of a good and updates all other priorities.
	//
	// It takes the id of the good to reprioritize and the new priority value.
	// It starts a database transaction to safely update all priorities.
	//
	// First it updates all goods with priority >= newPriority (except the prioritized good)
	// by incrementing their priority
	Reprioritize(ctx context.Context, id, newPriority int) (map[int]int, error)
}
