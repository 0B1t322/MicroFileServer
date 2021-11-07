package amqp

type Subscriber struct {
	Queue		string
	Consumer	string
	AutoAck		bool
	Exclusive	bool
	NoLocal		bool
	NoWait		bool
	Args		map[string]interface{}
}