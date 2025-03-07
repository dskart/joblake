# joblake

A lightweight and efficient worker pool implementation in Go. joblake makes it easy to process tasks concurrently with controlled parallelism.


## Installation

```bash
go get github.com/dskart/joblake
```

## Quick Start

```go
package main

import (
    "github.com/dskart/joblake"
)

func main() {
	numWokers := 5
	queueSize := 10
	jl := joblake.New(Config{
		NumWorkers: &numWokers,
		QueueSize:  &queueSize,
	})
    
	ctx := context.Background()
	logger := zerolog.New(os.Stdout)
	jl.Start(ctx, &logger)

	// Submit jobs
	for i := 0; i < 10; i++ {
		joblake.AddJob(Job{
			Name: "job-" + string(i),
			Fn: func(ctx context.Context) error {
				logger.Info().Msg("Running job")
				time.Sleep(1 * time.Second)
				logger.Info().Msg("Job done")
				return nil
			},
		})
	}
}
```
