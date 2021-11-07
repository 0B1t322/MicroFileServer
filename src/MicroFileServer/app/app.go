package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MicroFileServer/pkg/amqp/manager"
	"github.com/MicroFileServer/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

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
	Manager			manager.Manager
}

func New(cfg *config.Config) *App {
	app := &App{}
	app.Port = cfg.App.AppPort
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

	conn, err := amqp.Dial(cfg.AMQP.AMQPURI)
	if err != nil {
		log.WithFields(
			log.Fields{
				"package": "app",
				"func": "new",
				"err": err,
			},
		).Panic("failed to connect to RabbitMQ")
	}

	app.Manager, err = manager.NewManager(conn)
	if err != nil {
		log.WithFields(
			log.Fields{
				"package": "app",
				"func": "new",
				"err": err,
			},
		).Panic("failed to init rabbit Manager")
	}

	return app
}

func (a *App) AddApi(Builders ...apibuilder.ApiBulder) {
	for _, Builder := range Builders {
		Builder.AddLogger(a.Logger)
		Builder.AddAuthMiddleware(a.auth)
		Builder.CreateServices()
		Builder.Build(a.Router)
		Builder.BuildAMQP(a.Manager)
	}
}

func (a *App) Start() {
	log.Infof("Starting Application is port %s", a.Port)
	s := &http.Server{
		Addr: fmt.Sprintf(":%s",a.Port),
		Handler: a.Router,
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout: 2*time.Second,
	}

	go func() {
		log.Info("Staring AMQP Application layer")
		if err := a.Manager.Start(); err != nil {
			log.Panicf("Failed to start amqp application layer, err: %v", err)
		}
	}()

	if err := s.ListenAndServe(); err != nil {
		log.Panicf("Failed to start application %v", err)
	}
}