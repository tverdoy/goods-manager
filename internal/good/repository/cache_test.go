package repository

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goods-manager/internal/cache"
	"goods-manager/internal/cache/mocks"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	mocks2 "goods-manager/mocks"
	"testing"
)

func Test_goodRepositoryCache_Create(t *testing.T) {
	mockCache := mocks.NewCache(t)
	mockGoodRepo := mocks2.NewGoodRepository(t)

	cache := NewGoodRepositoryCache(mockCache, mockGoodRepo)

	good := entity.Good{
		Id:          523,
		ProjectId:   3,
		Name:        "Last",
		Description: "Go to home",
		Priority:    3,
		Removed:     false,
		CreatedAt:   "2024-03-05 12:00:00",
	}
	ctx := context.Background()

	mockGoodRepo.On("Create", ctx, &good).Return(nil)
	mockCache.On("Set", ctx, "good:523", &good).Return(nil)

	if err := cache.Create(ctx, &good); err != nil {
		t.Fatal(err)
	}
}

func Test_goodRepositoryCache_Delete(t *testing.T) {
	mockCache := mocks.NewCache(t)
	mockGoodRepo := mocks2.NewGoodRepository(t)

	cache := NewGoodRepositoryCache(mockCache, mockGoodRepo)

	good := entity.Good{
		Id:          523,
		ProjectId:   3,
		Name:        "Last",
		Description: "Go to home",
		Priority:    3,
		Removed:     false,
		CreatedAt:   "2024-03-05 12:00:00",
	}
	ctx := context.Background()

	mockGoodRepo.On("Delete", ctx, good.Id).Return(nil)
	mockCache.On("Remove", ctx, "good:523").Return(nil)

	if err := cache.Delete(ctx, good.Id); err != nil {
		t.Fatal(err)
	}
}

// # TODO add case if cache return nil
func Test_goodRepositoryCache_Get(t *testing.T) {
	mockCache := mocks.NewCache(t)
	mockGoodRepo := mocks2.NewGoodRepository(t)

	cache := NewGoodRepositoryCache(mockCache, mockGoodRepo)

	good := entity.Good{
		Id:          523,
		ProjectId:   3,
		Name:        "Last",
		Description: "Go to home",
		Priority:    3,
		Removed:     false,
		CreatedAt:   "2024-03-05 12:00:00",
	}
	ctx := context.Background()

	mockCache.On("Get", ctx, "good:523", &entity.Good{}).Return(nil).Run(func(args mock.Arguments) {
		mockGood := args.Get(2).(*entity.Good)
		*mockGood = good
	})

	goodCache, err := cache.Get(ctx, good.Id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &good, goodCache)
}

func Test_goodRepositoryCache_List(t *testing.T) {
	type fields struct {
		cache          cache.Cache
		goodRepository domain.GoodRepository
	}
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entity.Good
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goodRepositoryCache{
				cache:          tt.fields.cache,
				goodRepository: tt.fields.goodRepository,
			}
			got, err := g.List(tt.args.ctx, tt.args.limit, tt.args.offset)
			if !tt.wantErr(t, err, fmt.Sprintf("List(%v, %v, %v)", tt.args.ctx, tt.args.limit, tt.args.offset)) {
				return
			}
			assert.Equalf(t, tt.want, got, "List(%v, %v, %v)", tt.args.ctx, tt.args.limit, tt.args.offset)
		})
	}
}

func Test_goodRepositoryCache_Reprioritize(t *testing.T) {
	type fields struct {
		cache          cache.Cache
		goodRepository domain.GoodRepository
	}
	type args struct {
		ctx         context.Context
		id          int
		newPriority int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[int]int
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goodRepositoryCache{
				cache:          tt.fields.cache,
				goodRepository: tt.fields.goodRepository,
			}
			got, err := g.Reprioritize(tt.args.ctx, tt.args.id, tt.args.newPriority)
			if !tt.wantErr(t, err, fmt.Sprintf("Reprioritize(%v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.newPriority)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Reprioritize(%v, %v, %v)", tt.args.ctx, tt.args.id, tt.args.newPriority)
		})
	}
}

func Test_goodRepositoryCache_Update(t *testing.T) {
	type fields struct {
		cache          cache.Cache
		goodRepository domain.GoodRepository
	}
	type args struct {
		ctx  context.Context
		good *entity.Good
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &goodRepositoryCache{
				cache:          tt.fields.cache,
				goodRepository: tt.fields.goodRepository,
			}
			tt.wantErr(t, g.Update(tt.args.ctx, tt.args.good), fmt.Sprintf("Update(%v, %v)", tt.args.ctx, tt.args.good))
		})
	}
}
