package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MicroFileServer/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/MicroFileServer/pkg/apibuilder"
	"github.com/MicroFileServer/pkg/config"
	kl "github.com/go-kit/kit/log"
	klogrus "github.com/go-kit/kit/log/logrus"
	"github.com/MicroFileServer/service/middleware/auth"
)

type App struct {
	Router 			*mux.Router
	Repository		*repositories.Repositories
	Port			string
	Logger			kl.Logger
	auth			*auth.Auth

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
		log.Debug("Running in test mode")
		log.Debug("admin role",cfg.Auth.AdminRole)
		log.Debug("user role", cfg.Auth.UserRole)
		log.Debug("Audience", cfg.Auth.Audience)
	}
	app.auth = auth.NewAuth(
		&auth.Config{
			AuthConfig: cfg.Auth,
			Testmode: cfg.App.TestMode,
		},
	)

	app.Logger = klogrus.NewLogger(log.StandardLogger())

	return app
}

func (a *App) AddApi(Builders ...apibuilder.ApiBulder) {
	for _, Builder := range Builders {
		Builder.AddLogger(a.Logger)
		Builder.AddAuthMiddleware(a.auth)
		Builder.CreateServices()
		Builder.Build(a.Router)
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
	if err := s.ListenAndServe(); err != nil {
		log.Panicf("Failed to start application %v", err)
	}
}