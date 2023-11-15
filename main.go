// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/tinybit/go-plugin-log-example/shared"
)

func run() error {
	logWrapper, stderrToLogWriter := configureLogger()

	zlog.Info().Msg("Started stresshouse.")

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		Logger:           logWrapper,
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
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
	raw, err := rpcClient.Dispense("kv_grpc")
	if err != nil {
		return err
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	kv := raw.(shared.KV)
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
	// We don't want to see the plugin logs.
	// log.SetOutput(io.Discard)

	log.SetOutput(os.Stdout)
	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func configureLogger() (*LoggerWrapper, *StderrToLogWriter) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// Setup the default logger with global context
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampNano}
	logger := zerolog.New(output).With().Timestamp().Logger()
	// logger := zlog.With().Timestamp().Logger()
	loggerShed := logger.With().Str("app", "stresshouse").Logger()

	zlog.Logger = loggerShed

	baseLoggerOptions := hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Trace,
	}

	logWrapper := NewLoggerWrapper(baseLoggerOptions, &logger)
	stderrToLogWriter := NewStderrToLogWriter(&logger)

	return logWrapper, stderrToLogWriter
}
