package amqp

import (
	"github.com/MicroFileServer/pkg/amqp/manager"
	config "github.com/MicroFileServer/pkg/config/amqp"
	"github.com/streadway/amqp"

	amqptransport "github.com/go-kit/kit/transport/amqp"
)


type Consumer amqptransport.Subscriber

func (c Consumer) MakeHandler(ch *amqp.Channel) manager.Handler {
	return (amqptransport.Subscriber)(c).ServeDelivery(ch) 
}

type SubscriberWithConfig struct {
	Cfg		config.Subscriber
	*amqptransport.Subscriber
}

func (s *SubscriberWithConfig) MakeSubcriber() (
	config.Subscriber,
	manager.Consumer,
) {
	consumer := (*Consumer)(s.Subscriber)

	return s.Cfg, consumer
}