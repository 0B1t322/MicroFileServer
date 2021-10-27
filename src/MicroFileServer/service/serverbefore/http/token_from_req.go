package http

import (
	log "github.com/sirupsen/logrus"
	ctxtoken "github.com/MicroFileServer/pkg/contextvalue/token"
	"context"
	"net/http"
)

func TokenFromReq(
	ctx	context.Context,
	r	*http.Request,
) context.Context {
	log.WithFields(
		log.Fields{
			"package": "serverbefore/http",
			"func": "TokenFromReq",
		},
	).Debug("Token before")
	token := r.Header.Get("Authorization")
	if token == "" {
		return ctx
	}

	ctx = ctxtoken.New(
		ctx,
		token,
	)

	

	return ctx
}