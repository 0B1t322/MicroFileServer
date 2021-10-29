package apibuilder

import (
	"github.com/MicroFileServer/service/middleware/auth"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type ApiBulder interface {
	Build(*mux.Router)
	CreateServices()
	AddAuthMiddleware(auth.Auther)
	AddLogger(log.Logger)
	BuildGRPC(*grpc.Server)
}

