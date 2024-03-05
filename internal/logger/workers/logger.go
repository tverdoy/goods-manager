package workers

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	"goods-manager/internal/logger/usecase"
	"log"
	"sync"
	"time"
)

// LoggerWorker struct for receive and save logs
type LoggerWorker struct {
	nc            *nats.Conn
	loggerUsecase domain.LoggerUsecase
}

// Run start listing `usecase.Subject` and save logs to store
func (l *LoggerWorker) Run() error {
	mx := sync.Mutex{}
	buf := make([]*entity.Good, 0)

	go func() {
		for {
			mx.Lock()
			if len(buf) == 0 {
				mx.Unlock()
				continue
			}

			events := make([]*entity.Good, len(buf))
			copy(events, buf)
			buf = buf[:0]
			mx.Unlock()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			if err := l.loggerUsecase.SaveList(ctx, events); err != nil {
				log.Println("failed to save list of goods log", err)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	s, err := l.nc.Subscribe(usecase.Subject, func(m *nats.Msg) {
		// receive new message. Unmarshal and send to channel
		var good entity.Good
		if err := json.Unmarshal(m.Data, &good); err != nil {
			log.Println("failed unmarshal data from nats:", err)
		}

		mx.Lock()
		buf = append(buf, &good)
		mx.Unlock()
	})

	if s.IsValid() {
		log.Println("worker logger is valid")
	}

	return err
}

func NewLoggerWorker(nc *nats.Conn, loggerUsecase domain.LoggerUsecase) *LoggerWorker {
	return &LoggerWorker{nc: nc, loggerUsecase: loggerUsecase}
}
