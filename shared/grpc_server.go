// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	"github.com/tinybit/go-plugin-log-example/proto"
)

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	proto.UnimplementedKVServer
	Impl KV // This is the real implementation
}

func (m *GRPCServer) Ping(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Ping()
}

func (m *GRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Init(req.BrokerId)
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}
