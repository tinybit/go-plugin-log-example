// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/tinybit/go-plugin-log-example/proto"
	"google.golang.org/grpc"
)

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	proto.UnimplementedKVServer
	Impl          KV // This is the real implementation
	broker        *plugin.GRPCBroker
	brokerID      uint32
	logServerConn *grpc.ClientConn
	logClient     *GRPCLogHelperClient
}

func (m *GRPCServer) Ping(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Ping()
}

func (m *GRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
	m.brokerID = req.BrokerId

	err := m.connectToLoggerServer()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to logger server in stresshouse process: %v", err)
	}

	err = m.Impl.Init(m.brokerID)
	if err != nil {
		return nil, err
	}

	err = m.Impl.SetLogger(m.logClient)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}

func (m *GRPCServer) connectToLoggerServer() error {
	conn, err := m.broker.Dial(m.brokerID)
	if err != nil {
		return err
	}

	m.logServerConn = conn
	m.logClient = &GRPCLogHelperClient{proto.NewLogHelperClient(conn)}

	return nil
}

//defer conn.Close()
