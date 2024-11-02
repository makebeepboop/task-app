package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	taskgrpc "github.com/makebeepboop/protos/gen/go/task"
	"github.com/makebeepboop/task-app/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTaskAndStatus(t *testing.T) {
	ctx, st := suite.New(t)

	title := gofakeit.JobTitle()
	flow := gofakeit.Bird()
	number := gofakeit.Number(1, 1000)

	respCreate, err := st.TaskClient.Create(ctx, &taskgrpc.CreateRequest{
		Title:  title,
		Flow:   flow,
		Number: uint64(number),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respCreate.GetTaskId())

	respStatus, err := st.TaskClient.Status(ctx, &taskgrpc.StatusRequest{
		TaskId: respCreate.GetTaskId(),
	})
	require.NoError(t, err)
	assert.Equal(t, respStatus.GetStatus(), "in_progress")
}
