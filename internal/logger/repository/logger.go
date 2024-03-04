package repository

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
)

type loggerRepository struct {
	conn driver.Conn
}

func (l *loggerRepository) SaveList(ctx context.Context, goods []*entity.Good) error {
	query := `INSERT INTO goods (Id, ProjectId, Name, Description, Priority, Removed) VALUES `
	var values []interface{}

	for _, goods := range goods {
		query += "(?, ?, ?, ?, ?, ?),"
		values = append(values, goods.Id, goods.ProjectId, goods.Name, goods.Description, goods.Priority, goods.Removed)
	}

	return l.conn.AsyncInsert(ctx, query, false, values...)
}

func NewLoggerRepository(conn driver.Conn) domain.LoggerRepository {
	return &loggerRepository{conn: conn}
}
