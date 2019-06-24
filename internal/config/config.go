package config

import (
	"github.com/rs/zerolog"
	"go.soon.build/kit/config"
)

// application name
const APP_NAME = "bucket-boss"

// Config stores configuration options set by configuration file or env vars
type Config struct {
	Log Log
	AWS AWS
}

// Log contains logging configuration
type Log struct {
	Console bool
	Verbose bool
	Level   string
}

type AWS struct {
	Enabled    bool
	BucketName string
	Region     string
	AccessKey  string
	SecretKey  string
}

// Default is a default configuration setup with sane defaults
var Default = Config{
	Log{
		Level: zerolog.InfoLevel.String(),
	},
	AWS{},
}

// New constructs a new Config instance
func New(opts ...config.Option) (Config, error) {
	c := Default
	v := config.ViperWithDefaults("bucket-boss")
	err := config.ReadInConfig(v, &c, opts...)
	if err != nil {
		return c, err
	}
	return c, nil
}
