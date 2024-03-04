package repository

import (
	"context"
	"database/sql"
	"errors"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	"goods-manager/internal/transactor"
	"log"
)

type goodRepository struct {
	transactor *transactor.Transactor
}

// Create creates a new Good in the database.
//
// Parameters:
// - ctx: The context.Context for the operation.
// - good: A pointer to the Good entity to be created.
//
// Returns:
// - error: An error if the operation fails, nil otherwise.
//
// Description:
// The Create method inserts a new Good into the goods table in the database.
// It generates a new priority for the Good by incrementing the maximum priority of existing goods.
// The Good's project ID, name, description, and priority are provided as input parameters.
// The method returns the ID, priority, removed status, and creation timestamp of the newly
// created Good.
//
// Example usage:
//
//	good := &entity.Good{
//	    ProjectId:   1,
//	    Name:        "New Good",
//	    Description: "A new Good",
//	}
//	err := repository.Create(ctx, good)
//	if err != nil {
//	    log.Println("Failed to create Good:", err)
//	}
//
// Note:
// - The priority is calculated by incrementing the maximum priority of existing goods. If there are no existing goods, the priority will be set to 1.
// - The created_at field will be automatically set to the current timestamp by the database.
func (g *goodRepository) Create(ctx context.Context, good *entity.Good) error {
	query := `
        WITH max_priority AS (
            SELECT COALESCE(MAX(goods.priority), 0) AS priority FROM goods WHERE removed = false
        )
        
        INSERT INTO goods (project_id, name, description, priority) 
               VALUES ($1, $2, $3, (SELECT priority + 1 FROM max_priority))
        RETURNING id, priority, removed, created_at
    `

	tx, db := g.transactor.Connection(ctx)
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, good.ProjectId, good.Name, good.Description)
	} else {
		row = db.QueryRowContext(ctx, query, good.ProjectId, good.Name, good.Description)
	}

	return row.Scan(&good.Id, &good.Priority, &good.Removed, &good.CreatedAt)
}

// Get gets a Good from the database.
func (g *goodRepository) Get(ctx context.Context, id int) (*entity.Good, error) {
	query := `
		SELECT id, project_id, name, description, priority, removed, created_at FROM goods
			WHERE id = $1
	`

	tx, db := g.transactor.Connection(ctx)
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, id)
	} else {
		row = db.QueryRowContext(ctx, query, id)
	}

	var good entity.Good
	err := row.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrorGoodNotFound
		}

		return nil, err
	}

	return &good, nil
}

// Update updates a Good in the database.
func (g *goodRepository) Update(ctx context.Context, good *entity.Good) error {
	query := `
		UPDATE goods SET name = $1, description = $2 WHERE id = $3
	`

	tx, db := g.transactor.Connection(ctx)
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, good.Name, good.Description, good.Id)
	} else {
		_, err = db.ExecContext(ctx, query, good.Name, good.Description, good.Id)
	}

	return err
}

// Delete deletes a Good from the database.
func (g *goodRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE goods SET removed = true WHERE id = $1`

	tx, db := g.transactor.Connection(ctx)
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = db.ExecContext(ctx, query, id)
	}

	return err
}

// List gets a list of Goods from the database.
func (g *goodRepository) List(ctx context.Context, limit, offset int) ([]*entity.Good, error) {
	query := `
		SELECT id, project_id, name, description, priority, removed, created_at FROM goods
			LIMIT $1 OFFSET $2
	`

	tx, db := g.transactor.Connection(ctx)
	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, limit, offset)
	} else {
		rows, err = db.QueryContext(ctx, query, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panicln("failed closed rows", err)
		}
	}(rows)

	goods := make([]*entity.Good, 0)
	for rows.Next() {
		var good entity.Good
		err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			return nil, err
		}

		goods = append(goods, &good)
	}

	return goods, err
}

func (g *goodRepository) Reprioritize(ctx context.Context, id, newPriority int) (map[int]int, error) {
	queryUpdateAfter := `
		UPDATE goods SET priority = priority + 1 
		             WHERE priority >= $1 and id != $2
		             RETURNING id, priority;
	`

	tx, db := g.transactor.Connection(ctx)
	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.QueryContext(ctx, queryUpdateAfter, newPriority, id)
	} else {
		rows, err = db.QueryContext(ctx, queryUpdateAfter, newPriority, id)
	}

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panicln("failed closed rows", err)
		}
	}(rows)

	priorities := make(map[int]int)
	for rows.Next() {
		var id int
		var priority int
		err := rows.Scan(&id, &priority)
		if err != nil {
			return nil, err
		}

		priorities[id] = priority
	}

	queryUpdateGood := `UPDATE goods SET priority = $1 WHERE id = $2;`

	if tx != nil {
		_, err = tx.ExecContext(ctx, queryUpdateGood, newPriority, id)
	} else {
		_, err = db.ExecContext(ctx, queryUpdateGood, newPriority, id)
	}

	if err != nil {
		return nil, err
	}

	return priorities, nil
}

func NewGoodRepository(transactor *transactor.Transactor) domain.GoodRepository {
	return &goodRepository{transactor: transactor}
}
