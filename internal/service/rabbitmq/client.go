package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewClient(connection string) *MQClient {
	conn, err := amqp.Dial(connection)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return &MQClient{
		connection: conn,
		channel:    ch,
	}
}

func (c *MQClient) DeclareExchange(exchangeName string) {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")
}

func (c *MQClient) Send(data []byte, exchangeName string) {
	err := c.channel.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", data)
}

func (c *MQClient) Close() {
	c.connection.Close()
	c.channel.Close()
}
