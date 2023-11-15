// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.2
// source: kv.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	KV_Ping_FullMethodName = "/proto.KV/Ping"
	KV_Init_FullMethodName = "/proto.KV/Init"
	KV_Get_FullMethodName  = "/proto.KV/Get"
	KV_Put_FullMethodName  = "/proto.KV/Put"
)

// KVClient is the client API for KV service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*Empty, error)
}

type kVClient struct {
	cc grpc.ClientConnInterface
}

func NewKVClient(cc grpc.ClientConnInterface) KVClient {
	return &kVClient{cc}
}

func (c *kVClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, KV_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, KV_Init_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, KV_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, KV_Put_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVServer is the server API for KV service.
// All implementations must embed UnimplementedKVServer
// for forward compatibility
type KVServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	Init(context.Context, *InitRequest) (*Empty, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Put(context.Context, *PutRequest) (*Empty, error)
	mustEmbedUnimplementedKVServer()
}

// UnimplementedKVServer must be embedded to have forward compatible implementations.
type UnimplementedKVServer struct {
}

func (UnimplementedKVServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKVServer) Init(context.Context, *InitRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedKVServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKVServer) Put(context.Context, *PutRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedKVServer) mustEmbedUnimplementedKVServer() {}

// UnsafeKVServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVServer will
// result in compilation errors.
type UnsafeKVServer interface {
	mustEmbedUnimplementedKVServer()
}

func RegisterKVServer(s grpc.ServiceRegistrar, srv KVServer) {
	s.RegisterService(&KV_ServiceDesc, srv)
}

func _KV_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KV_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KV_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KV_Init_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KV_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KV_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KV_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KV_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KV_ServiceDesc is the grpc.ServiceDesc for KV service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KV_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.KV",
	HandlerType: (*KVServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KV_Ping_Handler,
		},
		{
			MethodName: "Init",
			Handler:    _KV_Init_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KV_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _KV_Put_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv.proto",
}

const (
	LogHelper_Log_FullMethodName = "/proto.LogHelper/Log"
)

// LogHelperClient is the client API for LogHelper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogHelperClient interface {
	Log(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*Empty, error)
}

type logHelperClient struct {
	cc grpc.ClientConnInterface
}

func NewLogHelperClient(cc grpc.ClientConnInterface) LogHelperClient {
	return &logHelperClient{cc}
}

func (c *logHelperClient) Log(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, LogHelper_Log_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogHelperServer is the server API for LogHelper service.
// All implementations must embed UnimplementedLogHelperServer
// for forward compatibility
type LogHelperServer interface {
	Log(context.Context, *LogRequest) (*Empty, error)
	mustEmbedUnimplementedLogHelperServer()
}

// UnimplementedLogHelperServer must be embedded to have forward compatible implementations.
type UnimplementedLogHelperServer struct {
}

func (UnimplementedLogHelperServer) Log(context.Context, *LogRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Log not implemented")
}
func (UnimplementedLogHelperServer) mustEmbedUnimplementedLogHelperServer() {}

// UnsafeLogHelperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogHelperServer will
// result in compilation errors.
type UnsafeLogHelperServer interface {
	mustEmbedUnimplementedLogHelperServer()
}

func RegisterLogHelperServer(s grpc.ServiceRegistrar, srv LogHelperServer) {
	s.RegisterService(&LogHelper_ServiceDesc, srv)
}

func _LogHelper_Log_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogHelperServer).Log(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogHelper_Log_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogHelperServer).Log(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LogHelper_ServiceDesc is the grpc.ServiceDesc for LogHelper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogHelper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LogHelper",
	HandlerType: (*LogHelperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Log",
			Handler:    _LogHelper_Log_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv.proto",
}
