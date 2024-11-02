package grpcapp

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	taskgrpc "github.com/makebeepboop/task-app/internal/transport/grpc/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

type Task interface {
	Create(
		ctx context.Context,
		title string,
		flow string,
		number uint64,
	) (taskId int64, err error)
	Status(
		ctx context.Context,
		taskId int64,
	) (isExist bool, status string, err error)
}

func New(
	log *slog.Logger,
	taskService Task,
	port int,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	taskgrpc.Register(gRPCServer, taskService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const operation = "grpcapp.Run"
	log := a.log.With(
		slog.String("operation", operation),
		slog.Int("port", a.port),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("gRPC server is running", slog.String("address", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (a *App) Stop() {
	const operation = "grpcapp.Stop"

	a.log.With(slog.String("operation", operation)).
		Info("stopping gRPC server")

	a.gRPCServer.Stop()
}
