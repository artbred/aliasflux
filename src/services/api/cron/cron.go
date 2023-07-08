package cron

import (
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/services/api/cron/jobs"
	"github.com/robfig/cron/v3"
	"time"
)

func runStartUpJobs(jobs ...cron.Job) {
	for _, job := range jobs {
		go job.Run()
	}
}

func Start() {
	c := cron.New(
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
		cron.WithLocation(time.UTC),
	)

	jobCache := jobs.NewCacheJob()

	_, err := c.AddJob("@every 1m", cron.NewChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	).Then(jobCache))

	if err != nil {
		common.Logger.WithError(err).Error("error adding job")
	}

	runStartUpJobs(jobCache)
	c.Start()
}
