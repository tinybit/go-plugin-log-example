// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared data between the host and plugins.
package shared

import (
	"github.com/hashicorp/go-plugin"
)

const (
	PluginProtocolVersion  = 1
	PluginMagicCookieKey   = "BASIC_PLUGIN"
	PluginMagicCookieValue = "hello"
	PluginID               = "kv_grpc"
)

func PluginHandshakeConfig() plugin.HandshakeConfig {
	// common handshake that is shared by plugin and host.
	var handshake = plugin.HandshakeConfig{
		// This isn't required when using VersionedPlugins
		ProtocolVersion:  PluginProtocolVersion,
		MagicCookieKey:   PluginMagicCookieKey,
		MagicCookieValue: PluginMagicCookieValue,
	}

	return handshake
}

func PluginMapClientConfig(pluginInstance plugin.Plugin) map[string]plugin.Plugin {
	// map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		PluginID: pluginInstance,
	}

	return pluginMap
}

func PluginMapServerConfig(kv KV) map[string]plugin.Plugin {
	// map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		PluginID: &KVGRPCPlugin{Impl: kv},
	}

	return pluginMap
}
