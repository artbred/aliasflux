package cron

import (
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
	_, _ = c.AddJob("@every 1m", jobCache)

	jobTld := jobs.NewTldJob()
	_, _ = c.AddJob("@every 24h", jobTld)

	runStartUpJobs(jobCache, jobTld)
	c.Start()
}
