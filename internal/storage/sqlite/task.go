package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/makebeepboop/task-app/internal/domain/models"
	"github.com/makebeepboop/task-app/internal/storage"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) CreateTask(
	ctx context.Context,
	title string,
	flow string,
	number uint64,
) (int64, error) {
	const operation = "storage.sqlite.CreateTask"

	stmt, err := s.db.Prepare("INSERT INTO tasks(title, flow, number, status) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	status := "in_progress"
	res, err := stmt.ExecContext(ctx, title, flow, number, status)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", operation, storage.ErrTaskAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return id, nil
}

func (s *Storage) Status(ctx context.Context, taskId int64) (models.Task, error) {
	const operation = "storage.sqlite.Task"

	stmt, err := s.db.Prepare("SELECT id, COALESCE(status, '') FROM tasks WHERE id = ?")
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, taskId)

	var task models.Task
	err = row.Scan(&task.ID, &task.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Task{}, fmt.Errorf("%s: %w", operation, storage.ErrTaskNotFound)
		}

		return models.Task{}, fmt.Errorf("%s: %w", operation, err)
	}

	return task, nil
}
