package manager

import (
	config "github.com/MicroFileServer/pkg/config/amqp"
	"github.com/streadway/amqp"
)

type Handler func(d *amqp.Delivery)


type Consumer interface{
	MakeHandler(ch *amqp.Channel) Handler
}

type Manager interface{
	// Start all consumers to listen amqp channel
	Start()	error


	AddConsumer(
		cfg			config.Subscriber,
		Consumer	Consumer,
	)
	
	// Init Queue in AMQP
	CreateQueue(
		cfg		config.Queue,
	) error
}

type manager struct {
	conn		*amqp.Connection
	ch			*amqp.Channel
	// Queues data
	// If create Queue data about them will be put here
	Queues		map[string]amqp.Queue

	Handlers	[]handler
}

func NewManager(
	conn	*amqp.Connection,
) (Manager, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &manager{
		conn: conn,
		ch: ch,
		Queues: make(map[string]amqp.Queue),
	}, nil
}

func (m *manager) CreateQueue(
	cfg		config.Queue,
) error {
	queue, err := m.ch.QueueDeclare(
		cfg.Name,
		cfg.Durable,
		cfg.AutoDelete,
		cfg.Exlusive,
		cfg.NoWait,
		cfg.Args,
	)
	if err != nil {
		return err
	}

	m.Queues[queue.Name] = queue

	return nil
}

func (m *manager) Start() error {
	for _, handler := range m.Handlers {
		if err := handler.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) AddConsumer(
	cfg			config.Subscriber,
	Consumer	Consumer,
) {
	m.Handlers = append(
		m.Handlers, 
		handler{
			Handler: Consumer.MakeHandler(m.ch),
			HandlerCreator: func() (<-chan amqp.Delivery, error) {
				return m.ch.Consume(
					cfg.Queue,
					cfg.Consumer,
					cfg.AutoAck,
					cfg.Exclusive,
					cfg.NoLocal,
					cfg.NoWait,
					cfg.Args,
				)
			},
		},
	)
}

type consumerCreater func() (<-chan amqp.Delivery, error)

type handler struct {
	Handler
	HandlerCreator		consumerCreater
	Messages			handlerMessages
}

type handlerMessages <-chan amqp.Delivery

func (h *handler) Start() error {
	messages, err := h.HandlerCreator()
	if err != nil {
		return err
	}

	h.Messages = messages

	go h.listenMessages()

	return nil
}

func (h *handler) listenMessages() {
	for message := range h.Messages {
		go h.Handler(&message)
	}
}
