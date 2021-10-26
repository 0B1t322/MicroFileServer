package http

import (
	"context"
	"net/http"

	"github.com/MicroFileServer/pkg/contextvalue/reqcontext"
)

func PutRequestInCTX(
	ctx	context.Context,
	r	*http.Request,
) context.Context {
	return reqcontext.New(
		ctx,
		r,
	)
}