package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"goods-manager/internal/domain/entity"
	"goods-manager/internal/transactor"
	"testing"
)

func initRepository() (*goodRepository, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	tr := transactor.NewTransactor(db)

	return &goodRepository{transactor: tr}, mock, nil
}

func Test_goodRepository_Create(t *testing.T) {
	repo, mock, err := initRepository()
	if err != nil {
		t.Fatal(err)
	}

	good := entity.Good{
		ProjectId:   4,
		Name:        "Good 1",
		Description: "Go to home",
	}

	oldGood := good

	createdAt := "2024-03-05 12:00:00"
	mock.ExpectQuery("INSERT INTO goods").
		WithArgs(good.ProjectId, good.Name, good.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "priority", "removed", "created_at"}).
			AddRow(1, 1, false, createdAt))

	if err := repo.Create(context.Background(), &good); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, 1, good.Id)
	assert.Equal(t, 1, good.Priority)
	assert.Equal(t, createdAt, good.CreatedAt)

	assert.Equal(t, oldGood.Name, good.Name)
	assert.Equal(t, oldGood.Description, good.Description)
}

func Test_goodRepository_Get(t *testing.T) {
	repo, mock, err := initRepository()
	if err != nil {
		t.Fatal(err)
	}

	good := &entity.Good{
		Id:          2,
		ProjectId:   4,
		Name:        "Good 1",
		Description: "Go to home",
		Priority:    3,
		Removed:     true,
		CreatedAt:   "2024-03-05 12:00:00",
	}

	mock.ExpectQuery("SELECT id, project_id, name, description, priority, removed, created_at FROM goods WHERE id = ?").
		WithArgs(good.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "name", "description", "priority", "removed", "created_at"}).
			AddRow(good.Id, good.ProjectId, good.Name, good.Description, good.Priority, good.Removed, good.CreatedAt))

	goodDb, err := repo.Get(context.Background(), good.Id)
	if err != nil {
		t.Errorf("Error getting good: %v", err)
	}

	assert.Equal(t, good, goodDb)
}

func Test_goodRepository_Update(t *testing.T) {
	repo, mock, err := initRepository()
	if err != nil {
		t.Fatal(err)
	}

	good := &entity.Good{
		Id:          2,
		ProjectId:   4,
		Name:        "Good 1",
		Description: "Go to home",
	}

	mock.ExpectExec("UPDATE goods").
		WithArgs(good.Name, good.Description, good.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), good)
	if err != nil {
		t.Errorf("Error updating good: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func Test_goodRepository_Delete(t *testing.T) {
	repo, mock, err := initRepository()
	if err != nil {
		t.Fatal(err)
	}

	id := 1
	mock.ExpectExec("UPDATE goods SET removed = true WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), id)

	if err != nil {
		t.Errorf("Error deleting good: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func Test_goodRepository_List(t *testing.T) {
	repo, mock, err := initRepository()
	if err != nil {
		t.Fatal(err)
	}

	limit := 10
	offset := 0

	mock.ExpectQuery("SELECT id, project_id, name, description, priority, removed, created_at FROM goods LIMIT ?").
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "name", "description", "priority", "removed", "created_at"}).
			AddRow(1, 1, "name_1", "description_1", 1, false, "2024-03-05 12:00:00").
			AddRow(2, 2, "name_2", "description_2", 2, false, "2024-03-06 12:00:00"))

	goods, err := repo.List(context.Background(), limit, offset)

	if err != nil {
		t.Errorf("Error listing goods: %v", err)
	}

	if len(goods) != 2 {
		t.Errorf("Expected 2 goods, got %d", len(goods))
	}

	if goods[0].Id != 1 || goods[0].ProjectId != 1 || goods[0].Name != "name_1" {
		t.Errorf("Unexpected content for the first good")
	}

	if goods[1].Id != 2 || goods[1].ProjectId != 2 || goods[1].Name != "name_2" {
		t.Errorf("Unexpected content for the second good")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func Test_goodRepository_Reprioritize(t *testing.T) {
	repo, mock, err := initRepository()
	if err != nil {
		t.Fatal(err)
	}

	id := 1
	newPriority := 5

	// Mock expected SQL query and its result for updating priorities
	mock.ExpectQuery("UPDATE goods").
		WithArgs(newPriority, id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "priority"}).
			AddRow(2, 6).
			AddRow(3, 7))

	// Mock expected SQL query and its result for updating the specified good's priority
	mock.ExpectExec("UPDATE goods").
		WithArgs(newPriority, id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the Reprioritize method
	priorities, err := repo.Reprioritize(context.Background(), id, newPriority)

	// Check if there was an error
	if err != nil {
		t.Errorf("Error reprioritizing goods: %v", err)
	}

	// Verify the returned priorities
	expectedPriorities := map[int]int{2: 6, 3: 7}
	for id, priority := range priorities {
		expectedPriority, ok := expectedPriorities[id]
		if !ok {
			t.Errorf("Unexpected id: %d returned in priorities", id)
		}
		if priority != expectedPriority {
			t.Errorf("Expected priority for id %d to be %d, got %d", id, expectedPriority, priority)
		}
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
