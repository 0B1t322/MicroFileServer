package amqp

type Queue struct {
	Name		string
	Durable		bool
	AutoDelete	bool
	Exlusive	bool
	NoWait		bool
	Args		map[string]interface{}
}