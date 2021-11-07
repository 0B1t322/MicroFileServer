package apibuilder

import (
	"github.com/MicroFileServer/pkg/amqp/manager"
	"github.com/MicroFileServer/service/middleware/auth"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type ApiBulder interface {
	Build(*mux.Router)
	BuildAMQP(manager.Manager)
	CreateServices()
	AddAuthMiddleware(auth.Auther)
	AddLogger(log.Logger)
}