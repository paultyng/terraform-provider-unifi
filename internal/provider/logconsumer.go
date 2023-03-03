package provider

import (
	"github.com/testcontainers/testcontainers-go"
)

type LogConsumer struct{}

func (c *LogConsumer) Accept(l testcontainers.Log) {
	testcontainers.Logger.Printf("[WARN] %s", string(l.Content))
}
