package entities

import (
	"github.com/murtll/mcserver-rcon/pkg/pb"
)

type DonateDelivery struct {
	*pb.DonateMessage

	DeliveryTag uint64
}
