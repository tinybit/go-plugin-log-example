// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

/*
	TODO:
	- add ability to specify log levels
	- add ability to specify key-value pairs in logs
*/

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
type KV struct {
	logClient shared.LogHelper
}

func NewKV() *KV {
	return &KV{}
}

func (k *KV) Ping() error {
	fmt.Fprintf(os.Stderr, "Plugin: got Ping() call.\n")
	return nil
}

func (k *KV) Init(uint32) error {
	fmt.Fprintf(os.Stderr, "Plugin: got Init() call.\n")

	return nil
}

func (k *KV) SetLogger(log shared.LogHelper) error {
	fmt.Fprintf(os.Stderr, "Plugin: got SetLogger() call.\n")

	k.logClient = log
	k.logClient.Log(0, "This is log message from Plugin.SetLogger()!")

	return nil
}

func (k *KV) Put(key string, value []byte) error {
	fmt.Fprintf(os.Stderr, "Plugin: got Put() call.\n")

	k.logClient.Log(0, "This is log message from Plugin.Put()!")

	value = []byte(fmt.Sprintf("value [%v] in plugin-go-grpc", string(value)))
	return os.WriteFile("kv_"+key, value, 0644)
}

func (k *KV) Get(key string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "Plugin: got Get() call.\n")

	k.logClient.Log(0, "This is log message from Plugin.Get()!")

	return os.ReadFile("kv_" + key)
}

func main() {
	serverInstance := NewKV()

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr,
		Level:  hclog.Debug,
	})

	plugin.Serve(&plugin.ServeConfig{
		Logger:          logger,
		HandshakeConfig: shared.PluginHandshakeConfig(),
		Plugins:         shared.PluginMapServerConfig(serverInstance),

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
