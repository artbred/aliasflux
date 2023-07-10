package jobs

import (
	"context"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/storages/redisdb"
	"github.com/sirupsen/logrus"
)

type Job struct {
	logger *logrus.Entry
}

func (j *Job) Run() {
	rdb := redisdb.Connection()
	rdb.FlushAll(context.Background())
}

func NewRedisJob() *Job {
	return &Job{
		logger: common.Logger.WithField("job", "redis"),
	}
}
