package main

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hovercross/natcat/pkg/reader"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	cfg, err := getConfig()
	flagOverride(&cfg)

	if err != nil {
		log.Fatalf("Unable to configure: %v", err)
	}

	nc, err := nats.Connect(cfg.Servers, cfg.options()...)
	if err != nil {
		log.Fatalf("Unable to connect to NATS: %v", err)
	}
	defer nc.Close()

	pubFunc := nc.Publish

	if cfg.Streaming {
		sc, err := stan.Connect(cfg.StreamClusterID, cfg.Name, stan.NatsConn(nc))

		if err != nil {
			log.Fatalf("Unable to connect ot NATS Streaming: %v", err)
		}

		pubFunc = sc.Publish
	}

	cfg.printConfig()

	log.Print("Connected to NATS")

	r := reader.Reader{
		Publish:           func(data []byte) error { return pubFunc(cfg.Topic, data) }, // Just injects the publication topic
		Input:             os.Stdin,
		Wrap:              cfg.Wrap,
		JSONInput:         cfg.WrapJSON,
		TimeGenerator:     time.Now,
		ReaderName:        cfg.Name,
		ReaderID:          uuid.New().String(),
		MessageIDFunction: func() string { return uuid.New().String() },
	}

	if err := r.Run(); err != nil {
		log.Fatalf("Unable to execute system: %v", err)
	}
}
