package files

import (
	"context"
	"fmt"
	"net/http"

	config "github.com/MicroFileServer/pkg/config/amqp"
	"github.com/MicroFileServer/pkg/statuscode"
	models "github.com/MicroFileServer/proto"
	adapter "github.com/MicroFileServer/service/adapters/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type AmqpServerConfig struct {
	Consumers
}

type Consumers struct {
	DeleteFile	config.Subscriber
}

type AmqpServer struct {
	Consumers	[]*adapter.SubscriberWithConfig
}

func NewAMQPServer(
	cfg			AmqpServerConfig,
	endpoints	Endpoints,
) *AmqpServer {
	a := &AmqpServer{
	}

	deleteFile := &adapter.SubscriberWithConfig{
		Subscriber: amqptransport.NewSubscriber(
			endpoints.DeleteFile,
			func(ctx context.Context, d *amqp.Delivery) (request interface{}, err error) {
				var deleteReq models.DeleteFileReq
				if err := proto.Unmarshal(d.Body, &deleteReq); err != nil {
					return nil, statuscode.WrapStatusError(
						fmt.Errorf("failed to decode proto file"),
						http.StatusBadRequest,
					)
				}
	
				return &DeleteFileReq{
					FileID: deleteReq.FileId,
				}, nil
			},
			func(_ context.Context, _ *amqp.Publishing, _ interface{}) error {
				return nil
			},
			amqptransport.SubscriberAfter(
				func(
					ctx context.Context, 
					d *amqp.Delivery, 
					ch amqptransport.Channel,
					p *amqp.Publishing,
				) context.Context {
					d.Ack(true)
					return ctx
				},
			),
		),
		Cfg: cfg.DeleteFile,
	}
	a.Consumers = append(a.Consumers, deleteFile)

	return a
}