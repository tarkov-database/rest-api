package health

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	cfg *config

	minDuration = 10 * time.Second
)

func init() {
	var err error

	cfg, err = newConfig()
	if err != nil {
		log.Printf("Configuration error: %s\n", err)
		os.Exit(2)
	}
}

type config struct {
	updateInterval   time.Duration
	latencyThreshold time.Duration
}

func newConfig() (*config, error) {
	c := &config{}

	if env := os.Getenv("HEALTHCHECK_INTERVAL"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			return c, fmt.Errorf("HEALTHCHECK_INTERVAL is not valid: %s", err)
		}

		if d < minDuration {
			c.updateInterval = minDuration
		} else {
			c.updateInterval = d
		}
	} else {
		c.updateInterval = 30 * time.Second
	}

	if env := os.Getenv("UNHEALTHY_LATENCY"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			return c, fmt.Errorf("UNHEALTHY_LATENCY is not valid: %s", err)
		}

		if d > c.updateInterval {
			return c, errors.New("UNHEALTHY_LATENCY can't be higher than HEALTHCHECK_INTERVAL")
		}

		c.latencyThreshold = d
	} else {
		c.latencyThreshold = 300 * time.Millisecond
	}

	return c, nil
}
