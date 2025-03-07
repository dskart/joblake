package joblake

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type JobLake struct {
	wp *workerPool
}

var (
	onceNewJL   sync.Once
	onceStartJL sync.Once
	jobLake     *JobLake
)

// New creates a new job lake
func New(cfg Config) *JobLake {
	numWorkers := 20
	queueSize := 1000

	if cfg.NumWorkers != nil {
		numWorkers = *cfg.NumWorkers
	}

	if cfg.QueueSize != nil {
		queueSize = *cfg.QueueSize
	}

	onceNewJL.Do(func() {
		wp := newWorkerPool(numWorkers, queueSize)
		jobLake = &JobLake{wp: wp}
	})
	return jobLake
}

type Job struct {
	Fn   func(ctx context.Context) error
	Name string
}

// Start starts the job lake
func (jl *JobLake) Start(ctx context.Context, logger *zerolog.Logger) {
	onceStartJL.Do(func() {
		jobLake.start(ctx, logger)
	})
}

func (jl *JobLake) start(ctx context.Context, logger *zerolog.Logger) {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return jl.wp.Run(ctx, logger)
	})
}

var ErrJobLakeNotInitialized = errors.New("Job lake is not initialized")

func AddJob(job Job) error {
	if jobLake == nil {
		return ErrJobLakeNotInitialized
	}
	jobLake.wp.AddJob(job)
	return nil
}
