package app

import (
	grpcapp "github.com/makebeepboop/task-app/internal/app/grpc"
	"github.com/makebeepboop/task-app/internal/services/task"
	"github.com/makebeepboop/task-app/internal/storage/sqlite"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := sqlite.New(
		log,
		storagePath,
	)
	if err != nil {
		panic(err)
	}

	taskService := task.New(
		log,
		storage,
		storage,
	)

	grpcApp := grpcapp.New(
		log,
		taskService,
		grpcPort,
	)

	return &App{
		GRPCServer: grpcApp,
	}
}
