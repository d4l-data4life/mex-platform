// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.0
// source: services/metadata/endpoints/blobs/blobs.proto

package pbBlobs

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
	Blobs_CreateBlob_FullMethodName = "/d4l.mex.blobs.Blobs/CreateBlob"
	Blobs_ListBlobs_FullMethodName  = "/d4l.mex.blobs.Blobs/ListBlobs"
	Blobs_GetBlob_FullMethodName    = "/d4l.mex.blobs.Blobs/GetBlob"
	Blobs_DeleteBlob_FullMethodName = "/d4l.mex.blobs.Blobs/DeleteBlob"
	Blobs_MeshTest_FullMethodName   = "/d4l.mex.blobs.Blobs/MeshTest"
)

// BlobsClient is the client API for Blobs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlobsClient interface {
	CreateBlob(ctx context.Context, in *CreateBlobRequest, opts ...grpc.CallOption) (*CreateBlobResponse, error)
	ListBlobs(ctx context.Context, in *ListBlobsRequest, opts ...grpc.CallOption) (*ListBlobsResponse, error)
	GetBlob(ctx context.Context, in *GetBlobRequest, opts ...grpc.CallOption) (*GetBlobResponse, error)
	DeleteBlob(ctx context.Context, in *DeleteBlobRequest, opts ...grpc.CallOption) (*DeleteBlobResponse, error)
	MeshTest(ctx context.Context, in *MeshTestRequest, opts ...grpc.CallOption) (*MeshTestResponse, error)
}

type blobsClient struct {
	cc grpc.ClientConnInterface
}

func NewBlobsClient(cc grpc.ClientConnInterface) BlobsClient {
	return &blobsClient{cc}
}

func (c *blobsClient) CreateBlob(ctx context.Context, in *CreateBlobRequest, opts ...grpc.CallOption) (*CreateBlobResponse, error) {
	out := new(CreateBlobResponse)
	err := c.cc.Invoke(ctx, Blobs_CreateBlob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobsClient) ListBlobs(ctx context.Context, in *ListBlobsRequest, opts ...grpc.CallOption) (*ListBlobsResponse, error) {
	out := new(ListBlobsResponse)
	err := c.cc.Invoke(ctx, Blobs_ListBlobs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobsClient) GetBlob(ctx context.Context, in *GetBlobRequest, opts ...grpc.CallOption) (*GetBlobResponse, error) {
	out := new(GetBlobResponse)
	err := c.cc.Invoke(ctx, Blobs_GetBlob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobsClient) DeleteBlob(ctx context.Context, in *DeleteBlobRequest, opts ...grpc.CallOption) (*DeleteBlobResponse, error) {
	out := new(DeleteBlobResponse)
	err := c.cc.Invoke(ctx, Blobs_DeleteBlob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobsClient) MeshTest(ctx context.Context, in *MeshTestRequest, opts ...grpc.CallOption) (*MeshTestResponse, error) {
	out := new(MeshTestResponse)
	err := c.cc.Invoke(ctx, Blobs_MeshTest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlobsServer is the server API for Blobs service.
// All implementations must embed UnimplementedBlobsServer
// for forward compatibility
type BlobsServer interface {
	CreateBlob(context.Context, *CreateBlobRequest) (*CreateBlobResponse, error)
	ListBlobs(context.Context, *ListBlobsRequest) (*ListBlobsResponse, error)
	GetBlob(context.Context, *GetBlobRequest) (*GetBlobResponse, error)
	DeleteBlob(context.Context, *DeleteBlobRequest) (*DeleteBlobResponse, error)
	MeshTest(context.Context, *MeshTestRequest) (*MeshTestResponse, error)
	mustEmbedUnimplementedBlobsServer()
}

// UnimplementedBlobsServer must be embedded to have forward compatible implementations.
type UnimplementedBlobsServer struct {
}

func (UnimplementedBlobsServer) CreateBlob(context.Context, *CreateBlobRequest) (*CreateBlobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBlob not implemented")
}
func (UnimplementedBlobsServer) ListBlobs(context.Context, *ListBlobsRequest) (*ListBlobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBlobs not implemented")
}
func (UnimplementedBlobsServer) GetBlob(context.Context, *GetBlobRequest) (*GetBlobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlob not implemented")
}
func (UnimplementedBlobsServer) DeleteBlob(context.Context, *DeleteBlobRequest) (*DeleteBlobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBlob not implemented")
}
func (UnimplementedBlobsServer) MeshTest(context.Context, *MeshTestRequest) (*MeshTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MeshTest not implemented")
}
func (UnimplementedBlobsServer) mustEmbedUnimplementedBlobsServer() {}

// UnsafeBlobsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlobsServer will
// result in compilation errors.
type UnsafeBlobsServer interface {
	mustEmbedUnimplementedBlobsServer()
}

func RegisterBlobsServer(s grpc.ServiceRegistrar, srv BlobsServer) {
	s.RegisterService(&Blobs_ServiceDesc, srv)
}

func _Blobs_CreateBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobsServer).CreateBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blobs_CreateBlob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobsServer).CreateBlob(ctx, req.(*CreateBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobs_ListBlobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBlobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobsServer).ListBlobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blobs_ListBlobs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobsServer).ListBlobs(ctx, req.(*ListBlobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobs_GetBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobsServer).GetBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blobs_GetBlob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobsServer).GetBlob(ctx, req.(*GetBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobs_DeleteBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobsServer).DeleteBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blobs_DeleteBlob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobsServer).DeleteBlob(ctx, req.(*DeleteBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobs_MeshTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MeshTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobsServer).MeshTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blobs_MeshTest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobsServer).MeshTest(ctx, req.(*MeshTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Blobs_ServiceDesc is the grpc.ServiceDesc for Blobs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Blobs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "d4l.mex.blobs.Blobs",
	HandlerType: (*BlobsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBlob",
			Handler:    _Blobs_CreateBlob_Handler,
		},
		{
			MethodName: "ListBlobs",
			Handler:    _Blobs_ListBlobs_Handler,
		},
		{
			MethodName: "GetBlob",
			Handler:    _Blobs_GetBlob_Handler,
		},
		{
			MethodName: "DeleteBlob",
			Handler:    _Blobs_DeleteBlob_Handler,
		},
		{
			MethodName: "MeshTest",
			Handler:    _Blobs_MeshTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/metadata/endpoints/blobs/blobs.proto",
}
