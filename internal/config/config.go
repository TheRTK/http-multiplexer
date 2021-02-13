package config

import "time"

type Config struct {
	PortHTTP                string
	RequestLimitCount       int
	MaxOutputRequestsPerURL int
	MaxInputRequestsCount   int
	OutputRequestTimeout    time.Duration
}
