package config

import (
	"fmt"

	"github.com/murtll/mcserver-rcon/pkg/util"
)

var Version = util.GetStrOrDefault("APP_VERSION", "0.1.0")

var ApiUrl = util.GetStrOrDefault("API_URL", "http://localhost:3001")
var ApiKey = util.GetStrOrDefault("API_KEY", "bigadminpassword")

var AmqpUrl = fmt.Sprintf("amqp://%s:%s@%s:%d",
	util.GetStrOrDefault("AMQP_USER", "guest"),
	util.GetStrOrDefault("AMQP_PASSWORD", "guest"),
	util.GetStrOrDefault("AMQP_HOST", "localhost"),
	util.GetIntOrDefault("AMQP_PORT", 5672))

var AmqpQueueName = util.GetStrOrDefault("AMQP_QUEUE_NAME", "donates")

var RconUrl = fmt.Sprintf("%s:%d", util.GetStrOrDefault("RCON_HOST", "localhost"), util.GetIntOrDefault("RCON_PORT", 25575))
var RconPass = util.GetStrOrDefault("RCON_PASS", "test")
