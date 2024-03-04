package usecase

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
)

const Subject = "logger:good"

type loggerUsecase struct {
	nc         *nats.Conn
	loggerRepo domain.LoggerRepository
}

// SendToQueue send message to `Subject`
func (l *loggerUsecase) SendToQueue(_ context.Context, good *entity.Good) error {
	data, err := json.Marshal(good)
	if err != nil {
		return err
	}

	return l.nc.Publish(Subject, data)
}

func (l *loggerUsecase) SaveList(ctx context.Context, goods []*entity.Good) error {
	return l.loggerRepo.SaveList(ctx, goods)
}

func NewLoggerUsecase(nc *nats.Conn, loggerRepo domain.LoggerRepository) domain.LoggerUsecase {
	return &loggerUsecase{nc: nc, loggerRepo: loggerRepo}
}
