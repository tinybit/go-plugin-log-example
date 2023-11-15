// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	zlog "github.com/rs/zerolog/log"
	"github.com/tinybit/go-plugin-log-example/proto"
)

type LogHelper interface {
	Log(level int, msg string) error
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCLogHelperClient struct{ client proto.LogHelperClient }

func (m *GRPCLogHelperClient) Log(level int, msg string) error {
	_, err := m.client.Log(context.Background(), &proto.LogRequest{
		Level:   int32(level),
		Message: msg,
	})

	if err != nil {
		zlog.Error().Msgf("Could not start log helper client: %v", err)
		return err
	}

	return nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCLogHelperServer struct {
	// This is the real implementation
	Impl LogHelper
}

func (m *GRPCLogHelperServer) Log(ctx context.Context, req *proto.LogRequest) (resp *proto.Empty, err error) {
	err = m.Impl.Log(int(req.GetLevel()), req.GetMessage())
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}
