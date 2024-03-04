package domain

import (
	"context"
	"goods-manager/internal/domain/entity"
)

type LoggerUsecase interface {
	SendToQueue(ctx context.Context, good *entity.Good) error
	SaveList(ctx context.Context, goods []*entity.Good) error
}

type LoggerRepository interface {
	SaveList(ctx context.Context, goods []*entity.Good) error
}
