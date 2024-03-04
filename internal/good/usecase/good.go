package usecase

import (
	"context"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	"goods-manager/internal/transactor"
)

// goodUsecase implementation `domain.GoodUsecase`.
//
// Mutable operation wrapper with transaction
// Some operation send log to queue.
type goodUsecase struct {
	goodRepo      domain.GoodRepository
	loggerUsecase domain.LoggerUsecase
	transactor    *transactor.Transactor
}

// Create new good and send log
func (g *goodUsecase) Create(ctx context.Context, good *entity.Good) error {
	return g.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		err := g.goodRepo.Create(ctx, good)

		if err == nil {
			if err := g.loggerUsecase.SendToQueue(ctx, good); err != nil {
				return err
			}
		}

		return err
	})
}

func (g *goodUsecase) Get(ctx context.Context, id int) (*entity.Good, error) {
	return g.goodRepo.Get(ctx, id)
}

// Update good and send log
func (g *goodUsecase) Update(ctx context.Context, good *entity.Good) error {
	return g.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		err := g.goodRepo.Update(ctx, good)

		if err == nil {
			if err := g.loggerUsecase.SendToQueue(ctx, good); err != nil {
				return err
			}
		}

		return err
	})
}

// Delete good and send log
func (g *goodUsecase) Delete(ctx context.Context, good *entity.Good) error {
	return g.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		good.Removed = true
		err := g.goodRepo.Delete(ctx, good.Id)

		if err == nil {
			if err := g.loggerUsecase.SendToQueue(ctx, good); err != nil {
				return err
			}
		}

		return err
	})
}

func (g *goodUsecase) List(ctx context.Context, limit, offset int) ([]*entity.Good, error) {
	return g.goodRepo.List(ctx, limit, offset)
}

func (g *goodUsecase) Reprioritize(ctx context.Context, id, newPriority int) (map[int]int, error) {
	var priorities map[int]int
	err := g.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		newPriorities, err := g.goodRepo.Reprioritize(ctx, id, newPriority)
		priorities = newPriorities
		return err
	})

	return priorities, err
}

func NewGoodUsecase(goodRepo domain.GoodRepository, loggerUsecase domain.LoggerUsecase, transactor *transactor.Transactor) domain.GoodUsecase {
	return &goodUsecase{goodRepo: goodRepo, loggerUsecase: loggerUsecase, transactor: transactor}
}
