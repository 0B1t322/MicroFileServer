package maxsize

import (
	"context"
	"errors"
)

var ErrRequestNotFound = errors.New("Request not found in ctx")

type maxSizeKey struct{}

type MaxSizeContext struct {
	context.Context
}

func New(
	ctx		context.Context,
	maxSize	int64,
) context.Context {
	return &MaxSizeContext{
		Context: context.WithValue(
			ctx,
			maxSizeKey{},
			maxSize,
		),
	}
}

func GetMaxSizeContext(
	ctx		context.Context,
) (int64, error) {
	if val, ok := ctx.Value(maxSizeKey{}).(int64); ok {
		return val, nil
	} else {
		return 0, ErrRequestNotFound
	}
}