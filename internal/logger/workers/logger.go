package workers

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	"goods-manager/internal/logger/usecase"
	"log"
	"time"
)

const BufferSize = 100 // Buffer size for send many logs to store

// LoggerWorker struct for receive and save logs
type LoggerWorker struct {
	nc            *nats.Conn
	loggerUsecase domain.LoggerUsecase
}

// Run start listing `usecase.Subject` and save logs to store
func (l *LoggerWorker) Run() error {
	ch := make(chan *entity.Good, 10)

	go func() {
		// goroutine for store events
		buf := [BufferSize]*entity.Good{}
		currentPtr := 0

		ctx := context.Background()

		for {
			// select good from chanel unless buffer size will
			// be `BufferSize` or end of chanel
		fillBuf:
			for {
				select {
				case good := <-ch:
					// add good to buffer
					buf[currentPtr] = good
					currentPtr++

					if currentPtr == BufferSize-1 {
						// buffer is full
						break fillBuf
					}
				default:
					break fillBuf
				}
			}

			// if buffer has goods, then store them
			if currentPtr > 0 {
				ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
				if err := l.loggerUsecase.SaveList(ctx, buf[:currentPtr]); err != nil {
					log.Println("failed to save list of goods log", err)
				}

				// clean buffer
				currentPtr = 0
				cancel()
			}

		}
	}()

	s, err := l.nc.Subscribe(usecase.Subject, func(m *nats.Msg) {
		// receive new message. Unmarshal and send to channel
		var good entity.Good
		if err := json.Unmarshal(m.Data, &good); err != nil {
			log.Println("failed unmarshal data from nats:", err)
		}

		ch <- &good
	})

	if s.IsValid() {
		log.Println("worker logger is valid")
	}

	return err
}

func NewLoggerWorker(nc *nats.Conn, loggerUsecase domain.LoggerUsecase) *LoggerWorker {
	return &LoggerWorker{nc: nc, loggerUsecase: loggerUsecase}
}
