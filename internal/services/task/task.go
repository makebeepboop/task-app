package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/makebeepboop/task-app/internal/domain/models"
	"github.com/makebeepboop/task-app/internal/storage"
	"github.com/makebeepboop/task-app/pkg/sl"
	"log/slog"
	"strconv"
)

type Task struct {
	log      *slog.Logger
	creater  Creater
	provider Provider
}

type Creater interface {
	CreateTask(
		ctx context.Context,
		title string,
		flow string,
		number uint64,
	) (taskId int64, err error)
}

type Provider interface {
	Status(ctx context.Context, taskId int64) (models.Task, error)
}

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
)

// New returns a new instance of the Task service.
func New(
	log *slog.Logger,
	creater Creater,
	provider Provider,
) *Task {
	return &Task{
		log:      log,
		creater:  creater,
		provider: provider,
	}
}

func (t *Task) Create(
	ctx context.Context,
	title string,
	flow string,
	number uint64,
) (int64, error) {
	const operation = "Task.Create"

	log := t.log.With(
		slog.String("operation", operation),
		slog.String("title", title),
		slog.String("flow", flow),
		slog.String("number", fmt.Sprint(number)),
	)

	log.Info("attempting to create task")

	taskId, err := t.creater.CreateTask(ctx, title, flow, number)
	if err != nil {
		if errors.Is(err, storage.ErrTaskAlreadyExists) {
			log.Warn("task already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", operation, err)
		}
		log.Error("failed to create task", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", operation, ErrTaskAlreadyExists)
	}

	log.Info("task created", slog.String("task_id", fmt.Sprint(taskId)))

	return taskId, nil
}

func (t *Task) Status(
	ctx context.Context,
	taskId int64,
) (bool, string, error) {
	const operation = "Task.Status"

	log := t.log.With(
		slog.String("operation", operation),
		slog.String("task_id", strconv.FormatInt(taskId, 10)),
	)

	log.Info("status check")

	task, err := t.provider.Status(ctx, taskId)
	if err != nil {
		log.Error("failed to get task", sl.Err(err))

		return false, "", fmt.Errorf("%s: %w", operation, err)
	}

	isExist := isTaskExist(task.ID)
	status := task.Status

	if !isExist {
		log.Info("task doesn't exist")
	} else {
		log.Info("task status", slog.String("status", status))
	}

	return isExist, status, nil
}

func isTaskExist(taskId int64) bool {
	return taskId > 0
}
