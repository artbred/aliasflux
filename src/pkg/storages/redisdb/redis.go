package redisdb

import (
	"context"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var connection *redis.Client
var once sync.Once

func Connection() *redis.Client {
	if connection == nil {
		once.Do(Init)
	}

	return connection
}

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ConnectionURLBuilder("redis"),
		Password: "",
		DB:       0,
	})

	status := rdb.Ping(context.Background())
	if status == nil {
		common.Logger.Error("Can't ping redis")
		time.Sleep(5 * time.Second)
		Init()
		return
	}

	if status.String() != "ping: PONG" {
		common.Logger.Error("Can't ping redis")
		time.Sleep(5 * time.Second)
		Init()
		return
	}

	connection = rdb
}

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}
}
