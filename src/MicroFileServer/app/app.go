package app

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/MicroFileServer/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/MicroFileServer/pkg/apibuilder"
	"github.com/MicroFileServer/pkg/config"
	"github.com/MicroFileServer/service/middleware/auth"
	kl "github.com/go-kit/kit/log"
	klogrus "github.com/go-kit/kit/log/logrus"
)

type App struct {
	Router 			*mux.Router
	Repository		*repositories.Repositories
	Port			string
	Logger			kl.Logger
	auth			*auth.Auth
	GRPCPort		string

	GRPCServer		*grpc.Server
}

func New(cfg *config.Config) *App {
	app := &App{}
	app.Port = cfg.App.AppPort
	app.GRPCPort = cfg.App.GrpcPort
	if _rep, err := repositories.New(&repositories.Config{
		DBURI: cfg.DB.URI,
	}); err != nil {
		log.WithFields(
			log.Fields{
				"package": "app",
				"func": "New",
				"err": err,
			},
		).Panic("Failed to init App")
	} else {
		app.Repository = _rep
	}

	app.Router = mux.NewRouter().PathPrefix("/api/mfs").Subrouter()
	if !cfg.App.TestMode {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
	app.auth = auth.NewAuth(
		&auth.Config{
			AuthConfig: cfg.Auth,
			Testmode: cfg.App.TestMode,
		},
	)

	app.Logger = klogrus.NewLogger(log.StandardLogger())
	app.GRPCServer = grpc.NewServer()

	return app
}

func (a *App) AddApi(Builders ...apibuilder.ApiBulder) {
	for _, Builder := range Builders {
		Builder.AddLogger(a.Logger)
		Builder.AddAuthMiddleware(a.auth)
		Builder.CreateServices()
		Builder.Build(a.Router)
		Builder.BuildGRPC(a.GRPCServer)
	}
}

func (a *App) Start() {
	if a.GRPCPort != "" {
		go a.StartGRPC()
	}
	a.StartHTTP()
}

func (a *App) StartHTTP() {
	log.Infof("Starting Application in port %s", a.Port)
	s := &http.Server{
		Addr: fmt.Sprintf(":%s",a.Port),
		Handler: a.Router,
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout: 2*time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Panicf("Failed to start application %v", err)
	}
}

func (a *App) StartGRPC() {
	log.Infof("Starting GRPC Application in port %s", a.GRPCPort)
	grpcListen, err := net.Listen("tcp", fmt.Sprintf(":%s", a.GRPCPort))
	if err != nil {
		log.Panicf("Failed to start grpc application %v", err)
	}

	if err := a.GRPCServer.Serve(grpcListen); err != nil {
		log.Panicf("Failed to start grpc application %v", err)
	}
}