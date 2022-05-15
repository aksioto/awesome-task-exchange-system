package rabbitmq

import (
	"log"
)

type RabbitmqService struct {
	server *MQServer
	client *MQClient
}

func NewRabbitmqService(connection string) *RabbitmqService {
	return &RabbitmqService{
		server: NewServer(connection),
		client: NewClient(connection),
	}
}

func (s *RabbitmqService) DeclareExchange(exchangeName string) {
	s.client.DeclareExchange(exchangeName)
	s.server.DeclareExchange(exchangeName)
}

func (s *RabbitmqService) Send(data string, exchangeName string) {
	s.client.Send(data, exchangeName)
}

func (s *RabbitmqService) Receive(callback Receiver, exchangeName string) {
	s.server.Receive(callback, exchangeName)
}

//TODO: todo
func (s *RabbitmqService) Close() {
	s.server.Close()
	s.client.Close()
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
