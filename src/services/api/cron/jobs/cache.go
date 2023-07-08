package jobs

import (
	"context"
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/storages/redisdb"
	"github.com/sirupsen/logrus"
)

type CacheJob struct {
	logger *logrus.Entry
}

func (j *CacheJob) cacheSettings() {
	ctx := context.Background()
	rdb := redisdb.Connection()

	settings, err := models.GetAllSettings()
	if err != nil {
		j.logger.WithError(err).Error("error getting settings")
		return
	}

	pgSettingsKeys := make(map[string]bool, len(settings))

	for _, setting := range settings {
		pgSettingsKeys[setting.Key] = true
		err = rdb.HSetNX(ctx, redisdb.SettingsKey.String(), setting.Key, setting.Value).Err()
		if err != nil {
			j.logger.WithError(err).Errorf("error setting %s", setting.Key)
		}
	}

	// Sync Redis with Postgres

	redisSettingsKeys, err := rdb.HKeys(ctx, redisdb.SettingsKey.String()).Result()
	if err != nil {
		j.logger.WithError(err).Error("error getting keys from Redis")
	}

	for _, key := range redisSettingsKeys {
		if !pgSettingsKeys[key] {
			err = rdb.HDel(ctx, redisdb.SettingsKey.String(), key).Err()
			if err != nil {
				j.logger.WithError(err).Errorf("error deleting %s", key)
			}
		}
	}
}

func (j *CacheJob) Run() {
	j.cacheSettings()
}

func NewCacheJob() *CacheJob {
	return &CacheJob{
		logger: common.Logger.WithField("job", "cache_job"),
	}
}
