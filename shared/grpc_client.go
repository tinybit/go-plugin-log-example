// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/tinybit/go-plugin-log-example/proto"
	"google.golang.org/grpc"
)

var (
	MainLogHelper LogHelper
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	ctx              context.Context
	broker           *plugin.GRPCBroker
	client           proto.KVClient
	brokerID         uint32
	isInitialized    bool
	logServerStarted bool
	mutex            sync.Mutex
}

func NewGRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) *GRPCClient {
	gClient := &GRPCClient{
		ctx:    ctx,
		broker: broker,
		client: proto.NewKVClient(conn),
	}

	return gClient
}

func (m *GRPCClient) Ping() error {
	_, err := m.client.Ping(m.ctx, &proto.Empty{})
	return err
}

func (m *GRPCClient) Initialize() error {
	return m.Init(m.brokerID)
}

func (m *GRPCClient) Init(brokerID uint32) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// initialize only once
	if m.isInitialized {
		return nil
	}

	_, err := m.client.Init(m.ctx, &proto.InitRequest{
		BrokerId: brokerID,
	})

	if err != nil {
		return err
	}

	m.isInitialized = true

	return nil
}

func (m *GRPCClient) Put(key string, value []byte) error {
	_, err := m.client.Put(m.ctx, &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *GRPCClient) Get(key string) ([]byte, error) {
	resp, err := m.client.Get(m.ctx, &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (m *GRPCClient) StartLogServer() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.logServerStarted {
		return
	}

	// start logger server and remember brokerID
	addHelperServer := &GRPCLogHelperServer{Impl: MainLogHelper}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterLogHelperServer(s, addHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	m.brokerID = brokerID
	m.logServerStarted = true
}
