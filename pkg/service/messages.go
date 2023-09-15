package service

import (
	"log"
	"strconv"
	"strings"

	"github.com/murtll/mcserver-rcon/pkg/entities"
	"github.com/murtll/mcserver-rcon/pkg/rcon"
	"github.com/murtll/mcserver-rcon/pkg/repository"
)

type MessageService struct {
	mr   *repository.MessageRepository
	ir   *repository.ItemRepository
	rcon *rcon.MCConn
}

func NewMessageService(mr *repository.MessageRepository, ir *repository.ItemRepository, rcon *rcon.MCConn) *MessageService {
	return &MessageService{
		mr:   mr,
		ir:   ir,
		rcon: rcon,
	}
}

func (ms *MessageService) Process() error {
	donates, err := ms.mr.ConsumeDonates()
	if err != nil {
		return err
	}

	for d := range donates {
		go func(donate *entities.DonateDelivery) {
			item, err := ms.ir.GetItem(int(donate.DonateItemId))
			if err != nil {
				log.Printf("unable to get item for donate '%v', error: '%v', rejecting...", donate, err)
				ms.mr.Reject(donate.DeliveryTag)
				return
			}

			cmd := strings.ReplaceAll(strings.ReplaceAll(item.Command, "%user%", donate.DonaterUsername), "%number%", strconv.Itoa(int(donate.Amount)))
			result, err := ms.rcon.SendCommand(cmd)
			if err != nil {
				log.Printf("unable to send rcon command '%s' for donate '%v', error: '%v', result '%v', rejecting...", cmd, donate, err, result)
				ms.mr.Reject(donate.DeliveryTag)
				return
			}
			log.Printf("sent '%v' to rcon, response: '%v'", cmd, result)
			ms.mr.Ack(donate.DeliveryTag)
		}(d)
	}
	return nil
}
