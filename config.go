package joblake

type Config struct {
	// You can pass in values for these fields to override the defaults
	NumWorkers *int `yaml:"NumWorkers" env:"NUM_WORKERS"`
	QueueSize  *int `yaml:"QueueSize" env:"QUEUE_SIZE"`
}
