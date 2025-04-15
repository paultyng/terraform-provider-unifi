// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.11.4
// source: github.com/moby/buildkit/session/secrets/secrets.proto

package secrets

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Secrets_GetSecret_FullMethodName = "/moby.buildkit.secrets.v1.Secrets/GetSecret"
)

// SecretsClient is the client API for Secrets service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretsClient interface {
	GetSecret(ctx context.Context, in *GetSecretRequest, opts ...grpc.CallOption) (*GetSecretResponse, error)
}

type secretsClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretsClient(cc grpc.ClientConnInterface) SecretsClient {
	return &secretsClient{cc}
}

func (c *secretsClient) GetSecret(ctx context.Context, in *GetSecretRequest, opts ...grpc.CallOption) (*GetSecretResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetSecretResponse)
	err := c.cc.Invoke(ctx, Secrets_GetSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretsServer is the server API for Secrets service.
// All implementations should embed UnimplementedSecretsServer
// for forward compatibility.
type SecretsServer interface {
	GetSecret(context.Context, *GetSecretRequest) (*GetSecretResponse, error)
}

// UnimplementedSecretsServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSecretsServer struct{}

func (UnimplementedSecretsServer) GetSecret(context.Context, *GetSecretRequest) (*GetSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecret not implemented")
}
func (UnimplementedSecretsServer) testEmbeddedByValue() {}

// UnsafeSecretsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretsServer will
// result in compilation errors.
type UnsafeSecretsServer interface {
	mustEmbedUnimplementedSecretsServer()
}

func RegisterSecretsServer(s grpc.ServiceRegistrar, srv SecretsServer) {
	// If the following call pancis, it indicates UnimplementedSecretsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Secrets_ServiceDesc, srv)
}

func _Secrets_GetSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServer).GetSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secrets_GetSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServer).GetSecret(ctx, req.(*GetSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Secrets_ServiceDesc is the grpc.ServiceDesc for Secrets service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Secrets_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "moby.buildkit.secrets.v1.Secrets",
	HandlerType: (*SecretsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSecret",
			Handler:    _Secrets_GetSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/moby/buildkit/session/secrets/secrets.proto",
}
