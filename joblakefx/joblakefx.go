package joblakefx

import (
	"context"

	"github.com/dskart/joblake"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Module("joblakefx",
	fx.Provide(New),
)

type Inputs struct {
	fx.In

	Lifecycle fx.Lifecycle

	Logger *zerolog.Logger
	Config joblake.Config `optional:"true"`
}

type Outputs struct {
	fx.Out

	JobLake *joblake.JobLake
}

func New(in Inputs) (Outputs, error) {
	jl := joblake.New(in.Config)
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			jl.Start(ctx, in.Logger)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return Outputs{
		JobLake: jl,
	}, nil
}
