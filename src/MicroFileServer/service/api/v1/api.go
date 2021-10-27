package v1

import (
	"github.com/MicroFileServer/pkg/repositories"
	"github.com/MicroFileServer/service/api/v1/files"
	"github.com/MicroFileServer/service/middleware/auth"
	"github.com/MicroFileServer/service/repoimp"
	kit_logger "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)


type Config struct {
	TestMode		bool
	MaxFileSizeMB	int64
}

type Api struct {
	Repo			*repositories.Repositories
	FileService		files.Service
	Auth			auth.Auther
	Logger			kit_logger.Logger
	TestMode		bool
	MaxFileSizeMB	int64
}

type ApiEndpoints struct {
	Files		files.Endpoints
}

func New(
	cfg		Config,
	Repo	*repositories.Repositories,
) *Api {
	return &Api{
		Repo: Repo,
		TestMode: cfg.TestMode,
		MaxFileSizeMB: cfg.MaxFileSizeMB,
	}
}

func (a *Api) AddLogger(logger kit_logger.Logger) {
	a.Logger = logger
}

func (a *Api) AddAuthMiddleware(auth auth.Auther) {
	a.Auth = auth
}

func (a *Api) CreateServices() {
	RepoImp := repoimp.New(a.Repo)

	a.FileService = files.New(
		RepoImp,
		a.Logger,
	)
}

func (a *Api) Build(r *mux.Router) {
	router := r.PathPrefix("/").Subrouter()
	var endpounts ApiEndpoints
	{
		if a.TestMode {
			endpounts = a._buildEndpoints()
		} else {
			endpounts = a.buildEndpoints()
		}
	}

	files.NewHTTPServer(
		&files.Config{
			MaxFileSizeMB: a.MaxFileSizeMB,
		},
		endpounts.Files,
		router,
	)
}