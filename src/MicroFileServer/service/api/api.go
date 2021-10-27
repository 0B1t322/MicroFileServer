package api

import (
	"github.com/MicroFileServer/pkg/repositories"
	"github.com/MicroFileServer/service/api/v1"
	"github.com/MicroFileServer/service/middleware/auth"
	kit_logger "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	swag "github.com/swaggo/http-swagger"
	_ "github.com/MicroFileServer/docs"
)

type Api struct {
	Repo		*repositories.Repositories
	Logger		kit_logger.Logger
	TestMode	bool
	V1			*v1.Api
}

type Config struct {
	V1Config	v1.Config
	
	// this field set testmode for all apis
	TestMode	bool
}

func New(
	cfg		Config,
	Repo	*repositories.Repositories,
) *Api {
	a := &Api{
		Repo: Repo,
		TestMode: cfg.TestMode,
	}

	if cfg.TestMode {
		cfg.V1Config.TestMode = true
	}

	a.V1 = v1.New(
		cfg.V1Config,
		Repo,
	)

	return a
}

func (a *Api) AddLogger(logger kit_logger.Logger) {
	a.V1.AddLogger(logger)
}

func (a *Api) AddAuthMiddleware(auth auth.Auther) {
	a.V1.AddAuthMiddleware(auth)
}

func (a *Api) CreateServices() {
	a.V1.CreateServices()
}

func (a *Api) Build(r *mux.Router) {
	a.V1.Build(r)

	docs := r.PathPrefix("/swagger")

	docs.Handler(
		swag.WrapHandler,
	)
}