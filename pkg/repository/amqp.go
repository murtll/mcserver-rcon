package repository

import (
	"log"

	"github.com/murtll/mcserver-rcon/pkg/entities"
	"github.com/murtll/mcserver-rcon/pkg/pb"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type MessageRepository struct {
	ch    *amqp.Channel
	queue *amqp.Queue
}

func NewMessageRepository(ch *amqp.Channel, qname string) (*MessageRepository, error) {
	q, err := ch.QueueDeclare(
		qname, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	return &MessageRepository{
		ch:    ch,
		queue: &q,
	}, nil
}

func (mr *MessageRepository) ConsumeDonates() (<-chan *entities.DonateDelivery, error) {
	msgs, err := mr.ch.Consume(
		mr.queue.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return nil, err
	}

	donates := make(chan *entities.DonateDelivery)

	go func() {
		defer mr.ch.Close()

		for m := range msgs {
			var donate *pb.DonateMessage
			err := proto.Unmarshal(m.Body, donate)
			if err != nil {
				log.Println("error unmarshalling message, rejecting...")
				m.Reject(false)
				continue
			}

			dd := &entities.DonateDelivery{
				DonateMessage: donate,
				DeliveryTag:   m.DeliveryTag,
			}

			donates <- dd
		}
	}()

	return donates, nil
}

func (mr *MessageRepository) Ack(tag uint64) {
	mr.ch.Ack(tag, false)
}

func (mr *MessageRepository) Reject(tag uint64) {
	mr.ch.Reject(tag, true)
}
