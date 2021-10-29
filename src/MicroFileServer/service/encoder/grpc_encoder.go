package encoder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/MicroFileServer/service/responce"
)

func EncodeGRPCResponce(
	ctx		context.Context,
	resp	interface{},
) (interface{}, error) {
	switch grpcresp := resp.(type) {
	case responce.GRPCResponce:
		return grpcresp.GRPCEncode(ctx)
	}
	return nil, statuscode.WrapStatusError(
		fmt.Errorf(""),
		http.StatusInternalServerError,
	)
}