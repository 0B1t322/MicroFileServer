package http

import (
	"context"
	"net/http"

	"github.com/MicroFileServer/pkg/contextvalue/maxsize"
	httptransport "github.com/go-kit/kit/transport/http"
)

func SetMaxBytesReader(
	maxSizeMB	int64,
) httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return maxsize.New(
			ctx,
			maxSizeMB * 1024 * 1024,
		)
	}
}