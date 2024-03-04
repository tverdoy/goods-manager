package repository

import (
	"context"
	"errors"
	"goods-manager/internal/cache"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	"strconv"
)

// goodRepositoryCache implementation `domain.GoodRepository`
// for proxying request by cache
type goodRepositoryCache struct {
	cache          cache.Cache
	goodRepository domain.GoodRepository
}

func (g *goodRepositoryCache) Create(ctx context.Context, good *entity.Good) error {
	if err := g.goodRepository.Create(ctx, good); err != nil {
		return err
	}

	return g.cache.Set(ctx, "good:"+strconv.Itoa(good.Id), good)
}

func (g *goodRepositoryCache) Get(ctx context.Context, id int) (*entity.Good, error) {
	good := &entity.Good{}
	err := g.cache.Get(ctx, "good:"+strconv.Itoa(id), good)

	if err != nil {
		if errors.Is(err, cache.ErrorNotExists) {
			// good not exists in cache, so get it from transactor
			good, err = g.goodRepository.Get(ctx, id)
			if err != nil {
				return nil, err
			}

			if err := g.cache.Set(ctx, "good:"+strconv.Itoa(id), good); err != nil {
				return nil, err
			}

			return good, nil
		}
		return nil, err
	}

	return good, nil
}

func (g *goodRepositoryCache) Update(ctx context.Context, good *entity.Good) error {
	if err := g.goodRepository.Update(ctx, good); err != nil {
		return err
	}

	return g.cache.Set(ctx, "good:"+strconv.Itoa(good.Id), good)
}

func (g *goodRepositoryCache) Delete(ctx context.Context, id int) error {
	if err := g.goodRepository.Delete(ctx, id); err != nil {
		return err
	}

	return g.cache.Remove(ctx, "good:"+strconv.Itoa(id))
}

func (g *goodRepositoryCache) List(ctx context.Context, limit, offset int) ([]*entity.Good, error) {
	return g.goodRepository.List(ctx, limit, offset)
}

func (g *goodRepositoryCache) Reprioritize(ctx context.Context, id, newPriority int) (map[int]int, error) {
	repositories, err := g.goodRepository.Reprioritize(ctx, id, newPriority)
	if err != nil {
		return nil, err
	}

	// Get good and update it priority
	for id, priority := range repositories {
		cacheKey := "good:" + strconv.Itoa(id)
		var good entity.Good
		if err := g.cache.Get(ctx, cacheKey, &good); err != nil {
			if errors.Is(err, cache.ErrorNotExists) {
				continue
			}

			return nil, err
		}

		good.Priority = priority
		if err := g.cache.Set(ctx, cacheKey, good); err != nil {
			return nil, err
		}
	}

	return repositories, nil
}

func NewGoodRepositoryCache(cache cache.Cache, goodRepository domain.GoodRepository) domain.GoodRepository {
	return &goodRepositoryCache{cache: cache, goodRepository: goodRepository}
}
