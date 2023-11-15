// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/tinybit/go-plugin-log-example/shared"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type KV struct{}

func (KV) Put(key string, value []byte) error {
	fmt.Fprintf(os.Stderr, "Plugin: got Put() call.\n")

	value = []byte(fmt.Sprintf("value [%v] in plugin-go-grpc", string(value)))
	return os.WriteFile("kv_"+key, value, 0644)
}

func (KV) Get(key string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "Plugin: got Get() call.\n")
	return os.ReadFile("kv_" + key)
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr,
		Level:  hclog.Debug,
	})

	plugin.Serve(&plugin.ServeConfig{
		Logger:          logger,
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv_grpc": &shared.KVGRPCPlugin{Impl: &KV{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
