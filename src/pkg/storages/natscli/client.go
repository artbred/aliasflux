package natscli

import (
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var connection *nats.Conn
var once sync.Once

func Connection() *nats.Conn {
	if connection == nil {
		once.Do(Init)
	}

	if connection == nil {
		logrus.Error("can't connect to nats")
		time.Sleep(1 * time.Second)
		Connection()
	}

	return connection
}

func Init() {
	nc, err := nats.Connect(config.ConnectionURLBuilder("nats"))
	if err != nil {
		logrus.Error(err)
		return
	}

	connection = nc
}
