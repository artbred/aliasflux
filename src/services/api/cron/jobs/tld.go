package jobs

import (
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/domain/providers/godaddy"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/storages/postgres"
	"github.com/sirupsen/logrus"
)

type TldJob struct {
	logger *logrus.Entry
}

func (j *TldJob) Run() {
	tldsGo, err := godaddy.NewClient().ListTld()
	if err != nil {
		j.logger.WithError(err).Error("failed to list tlds")
		return
	}

	conn := postgres.Connection()
	query := `INSERT INTO settings (platform, key, value)
		VALUES ($1, $2, $3::jsonb)
		ON CONFLICT (platform, key) DO UPDATE SET value = $3::jsonb`

	_, err = conn.Exec(query, models.PlatformDomain, models.SettingTld, tldsGo)
	if err != nil {
		j.logger.WithError(err).Error("failed to insert tlds")
		return
	}
}

func NewTldJob() *TldJob {
	return &TldJob{
		logger: common.Logger.WithField("job", "tld_job"),
	}
}
