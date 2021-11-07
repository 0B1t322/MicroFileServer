package main

import (
	"github.com/MicroFileServer/app"
	"github.com/MicroFileServer/pkg/config"
	"github.com/MicroFileServer/pkg/config/amqp"
	"github.com/MicroFileServer/service/api"
	v1 "github.com/MicroFileServer/service/api/v1"
	"github.com/MicroFileServer/service/api/v1/files"
)

// @title MicroFileService API
// @version 1.0
// @description This is a server for save and get files
// @BasePath /api/mfs
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg := config.GetConfig()
	app := app.New(cfg)
	app.AddApi(
		api.New(
			api.Config{
				TestMode: cfg.App.TestMode,
				V1Config: v1.Config{
					MaxFileSizeMB: cfg.App.MaxFileSize,
					AmqpLayerConfig: v1.AmqpLayerConfig{
						Files: files.AmqpServerConfig{
							Consumers: files.Consumers{
								DeleteFile: amqp.Subscriber{
									Queue: "/mfs/delete_file",
									Consumer: "mfs",
									AutoAck: false,
									NoLocal: false,
									NoWait: false,
									Exclusive: false,
									Args: nil,
								},
							},
						},
					},
				},
			},
			app.Repository,
		),
	)

	app.Start()
}