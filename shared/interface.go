// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/tinybit/go-plugin-log-example/proto"
)

type LogHelper interface {
	Log(level int, msg string) error
}

// KV is the interface that we're exposing as a plugin.
type KV interface {
	Ping() error
	Init(brokerID uint32) error
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type KVGRPCPlugin struct {
	plugin.Plugin
	Impl      KV
	ClientPtr *GRPCClient
}

func (p *KVGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterKVServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *KVGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	gClient := NewGRPCClient(ctx, broker, conn)
	p.ClientPtr = gClient

	return gClient, nil
}
