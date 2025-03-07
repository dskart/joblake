package joblake

import (
	"context"
	"sync/atomic"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type workerPool struct {
	numWorkers int
	queueSize  int
	jobQueue   chan Job
}

func newWorkerPool(numWorkers int, queueSize int) *workerPool {
	return &workerPool{
		numWorkers: numWorkers,
		queueSize:  queueSize,
		jobQueue:   make(chan Job, queueSize),
	}
}

var jobCounter atomic.Int32

func (p *workerPool) AddJob(job Job) {
	jobCounter.Add(1)
	p.jobQueue <- job
}

func (p *workerPool) Run(ctx context.Context, logger *zerolog.Logger) error {
	eg, ctx := errgroup.WithContext(ctx)
	for w := 1; w <= p.numWorkers; w++ {
		eg.Go(func() error {
			return worker(ctx, logger, p.jobQueue)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func worker(ctx context.Context, logger *zerolog.Logger, jobQueue <-chan Job) error {
	for {
		select {
		case Job := <-jobQueue:
			logger.Info().Str("jobName", Job.Name).Msg("starting job")
			err := Job.Fn(ctx)
			jobCounter.Add(-1)
			if err != nil {
				logger.Error().Str("jobName", Job.Name).Int32("jobsLeft", jobCounter.Load()).Err(err).Msg("job error")
			} else {
				logger.Info().Str("jobName", Job.Name).Int32("jobsLeft", jobCounter.Load()).Msg("job finished successfully")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
