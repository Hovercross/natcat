package main

import "flag"

// Flag can override the flags that were set by the config
func flagOverride(cfg *config) {
	flag.StringVar(&cfg.Servers, "servers", cfg.Servers, "Comma seperated list of servers. Overrides NATCAT_SERVERS.")
	flag.StringVar(&cfg.Topic, "topic", cfg.Topic, "Topic name to publish messages to. Overrides NATCAT_TOPIC.")
	flag.BoolVar(&cfg.Wrap, "wrap", cfg.Wrap, "Wrap input in an outer message. Overrides NATCAT_WRAP.")
	flag.BoolVar(&cfg.WrapJSON, "json", cfg.WrapJSON, "Specify input records as JSON, and don't double quote them. Overrides NATCAT_JSONINPUT.")
	flag.StringVar(&cfg.Name, "name", cfg.Name, "The name of the publisher. Overrides NATCAT_NAME.")
	flag.IntVar(&cfg.BufferSize, "buffersize", cfg.BufferSize, "The size of messages to buffer, if NATS goes offline. Overrides NATCAT_BUFFERSIZE.")
	flag.BoolVar(&cfg.Streaming, "stream", cfg.Streaming, "Use NATS Streaming instead of NATS native. Overrides NATCAT_STREAMING")
	flag.StringVar(&cfg.StreamClusterID, "streamClusterID", cfg.StreamClusterID, "Streaming Cluster ID for use with NATS Streaming: Overrides NATCAT_CLUSTERID")
	flag.Parse()
}
