// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	"github.com/hashicorp/go-plugin/examples/grpc/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	ctx    context.Context
	client proto.KVClient
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
