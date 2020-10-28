package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/nats-io/nats.go"
)

type config struct {
	Servers    string `env:"NATCAT_SERVERS" envDefault:"nats://localhost:4222"`
	Topic      string `env:"NATCAT_TOPIC" envDefault:"natcat"`
	Wrap       bool   `env:"NATCAT_WRAP" envDefault:"false"`
	WrapJSON   bool   `env:"NATCAT_JSONINPUT" envDefault:"false"`
	Name       string `env:"NATCAT_NAME" envDefault:"NatCat"`
	BufferSize int    `env:"NATCAT_BUFFERSIZE"`
}

func getConfig() (config, error) {
	cfg := config{}

	err := env.Parse(&cfg)

	return cfg, err
}

func (c config) options() []nats.Option {
	out := []nats.Option{}

	if c.Name != "" {
		out = append(out, nats.Name(c.Name))
	}

	if c.BufferSize > 0 {
		out = append(out, nats.ReconnectBufSize(c.BufferSize))
	}

	return out
}

func (c config) printConfig() {
	for _, server := range strings.Split(c.Servers, ",") {
		u, err := url.Parse(server)

		if err != nil {
			log.Printf("Unable to parse URL during display: %v", err)
			continue
		}

		log.Printf("Server: %s://%s\n", u.Scheme, u.Host)
	}

	log.Printf("Topic: %s\n", c.Topic)
	log.Printf("Wrap: %v\n", c.Wrap)
	log.Printf("JSON Input: %v\n", c.WrapJSON)
	log.Printf("Publisher name: %s\n", c.Name)
	log.Printf("Buffer size: %d", c.BufferSize)
}
