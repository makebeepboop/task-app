package taskgrpc

import (
	"context"
	taskpb "github.com/makebeepboop/protos/gen/go/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const emptyValue = 0

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

type serverAPI struct {
	taskpb.UnimplementedTaskServer
	task Task
}

func Register(gRPC *grpc.Server, task Task) {
	taskpb.RegisterTaskServer(gRPC, &serverAPI{task: task})
}

func (s *serverAPI) Create(
	ctx context.Context,
	request *taskpb.CreateRequest,
) (*taskpb.CreateResponse, error) {
	if err := validateCreateRequest(request); err != nil {
		return nil, err
	}

	taskId, err := s.task.Create(
		ctx,
		request.GetTitle(),
		request.GetFlow(),
		request.GetNumber(),
	)
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &taskpb.CreateResponse{
		TaskId: taskId,
	}, nil
}

func (s *serverAPI) Status(
	ctx context.Context,
	request *taskpb.StatusRequest,
) (*taskpb.StatusResponse, error) {
	if err := validateStatusRequest(request); err != nil {
		return nil, err
	}

	isExist, taskStatus, err := s.task.Status(ctx, request.GetTaskId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &taskpb.StatusResponse{
		IsExist: isExist,
		Status:  taskStatus,
	}, nil
}

func validateCreateRequest(request *taskpb.CreateRequest) error {
	if request.GetTitle() == "" {
		return status.Error(codes.InvalidArgument, "title is required")
	}

	if request.GetFlow() == "" {
		return status.Error(codes.InvalidArgument, "flow is required")
	}

	return nil
}

func validateStatusRequest(request *taskpb.StatusRequest) error {
	if request.GetTaskId() == emptyValue {
		return status.Error(codes.InvalidArgument, "task_id is required")
	}
	return nil
}
