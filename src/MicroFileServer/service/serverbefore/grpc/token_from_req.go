package grpc

import (
	"context"

	ctxtoken "github.com/MicroFileServer/pkg/contextvalue/token"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

func PutTokenIntoCTX(
	ctx		context.Context,
	md		metadata.MD,
) context.Context {
	log.Debug(md)
	tokenField := md.Get(":authority")
	var token string
	if len(tokenField) > 0 {
		token = tokenField[0]
	} else {
		token = ""
	}

	if token == "" {
		return ctx
	}

	ctx = ctxtoken.New(
		ctx,
		token,
	)

	return ctx
}