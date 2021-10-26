package subcontext

import (
	"context"
	"errors"
)

var (
	ErrSubNotFound = errors.New("Sub not found in context")
)

type SubKey struct{}

type SubContext struct {
	context.Context
}

func New(ctx context.Context, sub string) *SubContext {
	return &SubContext{
		Context: context.WithValue(
			ctx,
			SubKey{},
			sub,
		),
	}
}

func GetSubFromContext(
	ctx		context.Context,
) (string, error) {
	val := ctx.Value(SubKey{})

	if sub, ok := val.(string); !ok {
		return "", ErrSubNotFound
	} else {
		return sub, nil
	}
}