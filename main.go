// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/tinybit/go-plugin-log-example/shared"
)

const (
	MainProcessLogLabel   = "main"
	PluginProcessLogLabel = "plugin"
)

type LogHelper struct {
	logger *zerolog.Logger
}

func NewLogHelper(logger *zerolog.Logger) *LogHelper {
	lg := logger.With().Str("app", PluginProcessLogLabel).Logger()

	helper := &LogHelper{
		logger: &lg,
	}

	return helper
}

func (l *LogHelper) Log(level int, msg string) error {
	l.logger.Info().Msg(msg)
	return nil
}

func run() error {
	logger := configureLogger()
	shared.MainLogHelper = NewLogHelper(logger)
	logInjector := NewLogInjector(logger)
	stderrToLogWriter := NewStderrToLogWriter(logger)
	zlog.Logger = logger.With().Str("app", MainProcessLogLabel).Logger()

	zlog.Info().Msg("Started main process.")

	// We're a host. Start by launching the plugin process.
	pluginInstance := &shared.KVGRPCPlugin{}

	client := plugin.NewClient(&plugin.ClientConfig{
		Logger:           logInjector,
		HandshakeConfig:  shared.PluginHandshakeConfig(),
		Plugins:          shared.PluginMapClientConfig(pluginInstance),
		Cmd:              exec.Command("sh", "-c", os.Getenv("KV_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		SyncStderr:       stderrToLogWriter,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.PluginID)
	if err != nil {
		return err
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	kv := raw.(shared.KV)

	// ping first
	err = kv.Ping()
	if err != nil {
		return err
	}

	// init plugin
	err = pluginInstance.ClientPtr.Initialize()
	if err != nil {
		return err
	}

	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "get":
		result, err := kv.Get(os.Args[1])
		if err != nil {
			return err
		}

		zlog.Info().Msgf("Plugin get call result: %v", string(result))

	case "put":
		err := kv.Put(os.Args[1], []byte(os.Args[2]))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("please only use 'get' or 'put', given: %q", os.Args[0])
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func configureLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// Setup the default logger with global context
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	logger := zerolog.New(output).With().Timestamp().Logger()

	// logger := zlog.With().Timestamp().Logger()

	return &logger
}
