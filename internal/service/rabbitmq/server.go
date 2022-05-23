package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MQServer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      map[string]amqp.Queue
}

func NewServer(connection string) *MQServer {
	conn, err := amqp.Dial(connection)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return &MQServer{
		connection: conn,
		channel:    ch,
		queue:      map[string]amqp.Queue{},
	}
}

func (s *MQServer) DeclareExchange(exchangeName string) {
	err := s.channel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := s.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	s.queue[exchangeName] = q

	err = s.channel.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
}
func (s *MQServer) Receive(callback Receiver, exchangeName string) {
	msgs, err := s.channel.Consume(
		s.queue[exchangeName].Name, // queue
		"",                         // consumer
		true,                       // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			callback(d.Body)
			//log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*][%s] Waiting for messages.", exchangeName)
	<-forever
}

func (s *MQServer) Close() {
	s.connection.Close()
	s.channel.Close()
}
