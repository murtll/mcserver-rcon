package main

import (
	"log"

	"github.com/murtll/mcserver-rcon/pkg/config"
	"github.com/murtll/mcserver-rcon/pkg/rcon"
	"github.com/murtll/mcserver-rcon/pkg/repository"
	"github.com/murtll/mcserver-rcon/pkg/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ir, err := repository.NewItemRepository(config.ApiUrl, config.ApiKey)
	if err != nil {
		panic(err)
	}

	// open amqp connection
	conn, err := amqp.Dial(config.AmqpUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	mr, err := repository.NewMessageRepository(ch, config.AmqpQueueName)
	if err != nil {
		panic(err)
	}

	rcon, err := rcon.NewMCConn(config.RconUrl, config.RconPass)
	defer rcon.Close()

	ms := service.NewMessageService(mr, ir, rcon)

	log.Printf("Starting app v%s", config.Version)

	err = ms.Process()
	if err != nil {
		panic(err)
	}
}
