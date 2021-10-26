package encoder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/MicroFileServer/service/responce"
)

func EncodeResponce(
	ctx context.Context, 
	w http.ResponseWriter, 
	resp interface{},
) error {
	switch httpresp := resp.(type) {
	case responce.HTTPResponceCTX:
		httpresp.Headers(ctx, w)
		w.WriteHeader(httpresp.StatusCode())
		return httpresp.EncodeCTX(ctx, w)
	case responce.HTTPResponce:
		httpresp.Headers(ctx, w)
		w.WriteHeader(httpresp.StatusCode())
		return httpresp.Encode(w)
	}
	return statuscode.WrapStatusError(
		fmt.Errorf(""),
		http.StatusInternalServerError,
	)
}