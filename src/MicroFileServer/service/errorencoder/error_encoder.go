package errorencoder

import (
	"github.com/MicroFileServer/service/responce"
	"github.com/sirupsen/logrus"
	"net/http"
	"context"
)

func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	logrus.Debug("ServerErrorEncoder")
	resp := responce.FromErr(err)
	resp.WriteHeader(w)
	resp.WriteMessage(w)
}