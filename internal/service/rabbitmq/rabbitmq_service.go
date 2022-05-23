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
func (s *RabbitmqService) DeclareExchanges(exchangeNames []string) {
	for _, exchangeName := range exchangeNames {
		s.DeclareExchange(exchangeName)
	}
}

func (s *RabbitmqService) Send(data []byte, exchangeName string) error {
	s.client.Send(data, exchangeName)
	return nil
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
