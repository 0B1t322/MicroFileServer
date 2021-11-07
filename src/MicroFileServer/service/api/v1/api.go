package v1

import (
	"github.com/MicroFileServer/pkg/amqp/manager"
	"github.com/MicroFileServer/pkg/config/amqp"
	"github.com/MicroFileServer/pkg/repositories"
	"github.com/MicroFileServer/service/api/v1/files"
	"github.com/MicroFileServer/service/middleware/auth"
	"github.com/MicroFileServer/service/repoimp"
	kit_logger "github.com/go-kit/kit/log"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type AmqpLayerConfig struct {
	Files 	files.AmqpServerConfig
}

type Config struct {
	TestMode		bool
	MaxFileSizeMB	int64

	AmqpLayerConfig
}

type Api struct {
	Repo			*repositories.Repositories
	FileService		files.Service
	Auth			auth.Auther
	Logger			kit_logger.Logger
	TestMode		bool
	MaxFileSizeMB	int64

	Files			files.AmqpServerConfig
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
		Files: cfg.Files,
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

// for v1 apibuilder
func (a *Api) Build(r *mux.Router) {
	router := r.PathPrefix("/").Subrouter()
	endpoints := a.buildEndpoints()

	files.NewHTTPServer(
		&files.Config{
			MaxFileSizeMB: a.MaxFileSizeMB,
		},
		endpoints.Files,
		router,
	)
}



func (a *Api) BuildAMQP(Manager manager.Manager) {
	endpoints := a.buildAMQPEndpoints()

	if err := Manager.CreateQueue(
		amqp.Queue{
			Name: "/mfs/delete_file",
			Durable: true,
			AutoDelete: true,
			Exlusive: false,
			NoWait: false,
		},
	); err != nil {
		logrus.Errorf("Failed to create Queue: %v", err)
	}

	filesServer := files.NewAMQPServer(
		a.Files,
		endpoints.Files,
	)

	for _, consumer := range filesServer.Consumers {
		Manager.AddConsumer(
			consumer.MakeSubcriber(),
		)
	}
}