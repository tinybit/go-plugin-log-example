// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	zlog "github.com/rs/zerolog/log"
	"github.com/tinybit/go-plugin-log-example/proto"
	"google.golang.org/grpc"
)

var (
	MainLogHelper LogHelper
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	ctx           context.Context
	broker        *plugin.GRPCBroker
	client        proto.KVClient
	isInitialized bool
	mutex         sync.Mutex
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
	return m.Init(0)
}

func (m *GRPCClient) Init(uint32) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// initialize only once
	if m.isInitialized {
		return nil
	}

	zlog.Info().Msg("Starting logger server...")
	brokerID := m.startLogServer(MainLogHelper)

	fmt.Println("GRPCClient.Init, brokerID:", brokerID)

	_, err := m.client.Init(m.ctx, &proto.InitRequest{
		BrokerId: brokerID,
	})

	if err != nil {
		return err
	}

	m.isInitialized = true

	return nil
}

func (m *GRPCClient) SetLogger(LogHelper) error {
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

func (m *GRPCClient) startLogServer(log LogHelper) (brokerID uint32) {
	// start logger server and remember brokerID
	addHelperServer := &GRPCLogHelperServer{Impl: log}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterLogHelperServer(s, addHelperServer)

		return s
	}

	brokerID = m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	return
}
