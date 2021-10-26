package reqcontext

import (
	"context"
	"errors"
	"net/http"
)

var ErrRequestNotFound = errors.New("Request not found in ctx")

type requestKey struct{}

type RequestContext struct {
	context.Context
}

func New(
	ctx		context.Context,
	req		*http.Request,
) context.Context {
	return &RequestContext{
		Context: context.WithValue(
			ctx,
			requestKey{},
			req,
		),
	}
}

func GetRequestFromContext(
	ctx		context.Context,
) (*http.Request, error) {
	if val, ok := ctx.Value(requestKey{}).(*http.Request); ok {
		return val, nil
	} else {
		return nil, ErrRequestNotFound
	}
}