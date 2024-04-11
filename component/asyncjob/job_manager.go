package asyncjob

import (
	"context"
	"log"
)

type group struct {
	jobs         []Job
	isConcurrent bool
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	return &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
	}
}

func (g *group) Run(ctx context.Context) error {
	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				errChan <- g.runJob(ctx, aj)
			}(g.jobs[i])
			continue
		}

		job := g.jobs[i]
		errChan <- g.runJob(ctx, job)
	}

	for i := 1; i <= len(g.jobs); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// Retry if needed
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}
