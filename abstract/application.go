package abstract

import "context"

type (
	Operation[TReq any, TRes any] interface {
		Execute(ctx context.Context, req TReq) (TRes, error)
	}
)
