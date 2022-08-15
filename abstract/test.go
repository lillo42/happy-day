package abstract

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockOperation[TReq any, TRes any] struct {
	mock.Mock
}

func (mock *MockOperation[TReq, TRes]) Execute(ctx context.Context, req TReq) (TRes, error) {
	args := mock.Called(ctx, req)
	return args.Get(0).(TRes), args.Error(1)
}
