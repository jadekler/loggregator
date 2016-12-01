// Code generated by protoc-gen-go.
// source: grpc.proto
// DO NOT EDIT!

/*
Package plumbing is a generated protocol buffer package.

It is generated from these files:
	grpc.proto

It has these top-level messages:
	EnvelopeData
	PushResponse
	SubscriptionRequest
	Filter
	Response
	ContainerMetricsRequest
	ContainerMetricsResponse
	RecentLogsRequest
	RecentLogsResponse
*/
package plumbing

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EnvelopeData struct {
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *EnvelopeData) Reset()                    { *m = EnvelopeData{} }
func (m *EnvelopeData) String() string            { return proto.CompactTextString(m) }
func (*EnvelopeData) ProtoMessage()               {}
func (*EnvelopeData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type PushResponse struct {
}

func (m *PushResponse) Reset()                    { *m = PushResponse{} }
func (m *PushResponse) String() string            { return proto.CompactTextString(m) }
func (*PushResponse) ProtoMessage()               {}
func (*PushResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type SubscriptionRequest struct {
	ShardID string  `protobuf:"bytes,1,opt,name=shardID" json:"shardID,omitempty"`
	Filter  *Filter `protobuf:"bytes,2,opt,name=filter" json:"filter,omitempty"`
}

func (m *SubscriptionRequest) Reset()                    { *m = SubscriptionRequest{} }
func (m *SubscriptionRequest) String() string            { return proto.CompactTextString(m) }
func (*SubscriptionRequest) ProtoMessage()               {}
func (*SubscriptionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SubscriptionRequest) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type Filter struct {
	AppID string `protobuf:"bytes,1,opt,name=appID" json:"appID,omitempty"`
}

func (m *Filter) Reset()                    { *m = Filter{} }
func (m *Filter) String() string            { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()               {}
func (*Filter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// Note: Ideally this would be EnvelopeData but for the time being we do not
// want to pay the cost of planning an upgrade path for this to be renamed.
type Response struct {
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type ContainerMetricsRequest struct {
	AppID string `protobuf:"bytes,1,opt,name=appID" json:"appID,omitempty"`
}

func (m *ContainerMetricsRequest) Reset()                    { *m = ContainerMetricsRequest{} }
func (m *ContainerMetricsRequest) String() string            { return proto.CompactTextString(m) }
func (*ContainerMetricsRequest) ProtoMessage()               {}
func (*ContainerMetricsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type ContainerMetricsResponse struct {
	Payload [][]byte `protobuf:"bytes,1,rep,name=payload,proto3" json:"payload,omitempty"`
}

func (m *ContainerMetricsResponse) Reset()                    { *m = ContainerMetricsResponse{} }
func (m *ContainerMetricsResponse) String() string            { return proto.CompactTextString(m) }
func (*ContainerMetricsResponse) ProtoMessage()               {}
func (*ContainerMetricsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type RecentLogsRequest struct {
	AppID string `protobuf:"bytes,1,opt,name=appID" json:"appID,omitempty"`
}

func (m *RecentLogsRequest) Reset()                    { *m = RecentLogsRequest{} }
func (m *RecentLogsRequest) String() string            { return proto.CompactTextString(m) }
func (*RecentLogsRequest) ProtoMessage()               {}
func (*RecentLogsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type RecentLogsResponse struct {
	Payload [][]byte `protobuf:"bytes,1,rep,name=payload,proto3" json:"payload,omitempty"`
}

func (m *RecentLogsResponse) Reset()                    { *m = RecentLogsResponse{} }
func (m *RecentLogsResponse) String() string            { return proto.CompactTextString(m) }
func (*RecentLogsResponse) ProtoMessage()               {}
func (*RecentLogsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func init() {
	proto.RegisterType((*EnvelopeData)(nil), "plumbing.EnvelopeData")
	proto.RegisterType((*PushResponse)(nil), "plumbing.PushResponse")
	proto.RegisterType((*SubscriptionRequest)(nil), "plumbing.SubscriptionRequest")
	proto.RegisterType((*Filter)(nil), "plumbing.Filter")
	proto.RegisterType((*Response)(nil), "plumbing.Response")
	proto.RegisterType((*ContainerMetricsRequest)(nil), "plumbing.ContainerMetricsRequest")
	proto.RegisterType((*ContainerMetricsResponse)(nil), "plumbing.ContainerMetricsResponse")
	proto.RegisterType((*RecentLogsRequest)(nil), "plumbing.RecentLogsRequest")
	proto.RegisterType((*RecentLogsResponse)(nil), "plumbing.RecentLogsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Doppler service

type DopplerClient interface {
	Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Doppler_SubscribeClient, error)
	ContainerMetrics(ctx context.Context, in *ContainerMetricsRequest, opts ...grpc.CallOption) (*ContainerMetricsResponse, error)
	RecentLogs(ctx context.Context, in *RecentLogsRequest, opts ...grpc.CallOption) (*RecentLogsResponse, error)
}

type dopplerClient struct {
	cc *grpc.ClientConn
}

func NewDopplerClient(cc *grpc.ClientConn) DopplerClient {
	return &dopplerClient{cc}
}

func (c *dopplerClient) Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Doppler_SubscribeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Doppler_serviceDesc.Streams[0], c.cc, "/plumbing.Doppler/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &dopplerSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Doppler_SubscribeClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type dopplerSubscribeClient struct {
	grpc.ClientStream
}

func (x *dopplerSubscribeClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dopplerClient) ContainerMetrics(ctx context.Context, in *ContainerMetricsRequest, opts ...grpc.CallOption) (*ContainerMetricsResponse, error) {
	out := new(ContainerMetricsResponse)
	err := grpc.Invoke(ctx, "/plumbing.Doppler/ContainerMetrics", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dopplerClient) RecentLogs(ctx context.Context, in *RecentLogsRequest, opts ...grpc.CallOption) (*RecentLogsResponse, error) {
	out := new(RecentLogsResponse)
	err := grpc.Invoke(ctx, "/plumbing.Doppler/RecentLogs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Doppler service

type DopplerServer interface {
	Subscribe(*SubscriptionRequest, Doppler_SubscribeServer) error
	ContainerMetrics(context.Context, *ContainerMetricsRequest) (*ContainerMetricsResponse, error)
	RecentLogs(context.Context, *RecentLogsRequest) (*RecentLogsResponse, error)
}

func RegisterDopplerServer(s *grpc.Server, srv DopplerServer) {
	s.RegisterService(&_Doppler_serviceDesc, srv)
}

func _Doppler_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DopplerServer).Subscribe(m, &dopplerSubscribeServer{stream})
}

type Doppler_SubscribeServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type dopplerSubscribeServer struct {
	grpc.ServerStream
}

func (x *dopplerSubscribeServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Doppler_ContainerMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContainerMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DopplerServer).ContainerMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plumbing.Doppler/ContainerMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DopplerServer).ContainerMetrics(ctx, req.(*ContainerMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Doppler_RecentLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecentLogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DopplerServer).RecentLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plumbing.Doppler/RecentLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DopplerServer).RecentLogs(ctx, req.(*RecentLogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Doppler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "plumbing.Doppler",
	HandlerType: (*DopplerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ContainerMetrics",
			Handler:    _Doppler_ContainerMetrics_Handler,
		},
		{
			MethodName: "RecentLogs",
			Handler:    _Doppler_RecentLogs_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Doppler_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

// Client API for DopplerIngestor service

type DopplerIngestorClient interface {
	Pusher(ctx context.Context, opts ...grpc.CallOption) (DopplerIngestor_PusherClient, error)
}

type dopplerIngestorClient struct {
	cc *grpc.ClientConn
}

func NewDopplerIngestorClient(cc *grpc.ClientConn) DopplerIngestorClient {
	return &dopplerIngestorClient{cc}
}

func (c *dopplerIngestorClient) Pusher(ctx context.Context, opts ...grpc.CallOption) (DopplerIngestor_PusherClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DopplerIngestor_serviceDesc.Streams[0], c.cc, "/plumbing.DopplerIngestor/Pusher", opts...)
	if err != nil {
		return nil, err
	}
	x := &dopplerIngestorPusherClient{stream}
	return x, nil
}

type DopplerIngestor_PusherClient interface {
	Send(*EnvelopeData) error
	CloseAndRecv() (*PushResponse, error)
	grpc.ClientStream
}

type dopplerIngestorPusherClient struct {
	grpc.ClientStream
}

func (x *dopplerIngestorPusherClient) Send(m *EnvelopeData) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dopplerIngestorPusherClient) CloseAndRecv() (*PushResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PushResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DopplerIngestor service

type DopplerIngestorServer interface {
	Pusher(DopplerIngestor_PusherServer) error
}

func RegisterDopplerIngestorServer(s *grpc.Server, srv DopplerIngestorServer) {
	s.RegisterService(&_DopplerIngestor_serviceDesc, srv)
}

func _DopplerIngestor_Pusher_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DopplerIngestorServer).Pusher(&dopplerIngestorPusherServer{stream})
}

type DopplerIngestor_PusherServer interface {
	SendAndClose(*PushResponse) error
	Recv() (*EnvelopeData, error)
	grpc.ServerStream
}

type dopplerIngestorPusherServer struct {
	grpc.ServerStream
}

func (x *dopplerIngestorPusherServer) SendAndClose(m *PushResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dopplerIngestorPusherServer) Recv() (*EnvelopeData, error) {
	m := new(EnvelopeData)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _DopplerIngestor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "plumbing.DopplerIngestor",
	HandlerType: (*DopplerIngestorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Pusher",
			Handler:       _DopplerIngestor_Pusher_Handler,
			ClientStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("grpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x52, 0x4d, 0x4f, 0x83, 0x40,
	0x10, 0x2d, 0x1a, 0x69, 0x3b, 0x36, 0x5a, 0x47, 0xa3, 0x04, 0x3f, 0x52, 0x89, 0x07, 0xbc, 0xa0,
	0xa9, 0x1e, 0x3d, 0x69, 0x35, 0x69, 0xa2, 0xd1, 0xe0, 0xc9, 0x78, 0x02, 0x3a, 0x52, 0x12, 0xdc,
	0x5d, 0x77, 0x17, 0x13, 0x7f, 0xb8, 0x77, 0xd3, 0x0f, 0x0a, 0x56, 0x6c, 0x8f, 0xf3, 0xf5, 0xe6,
	0xcd, 0x9b, 0x07, 0x10, 0x4b, 0x11, 0x79, 0x42, 0x72, 0xcd, 0xb1, 0x21, 0xd2, 0xec, 0x3d, 0x4c,
	0x58, 0xec, 0xb8, 0xd0, 0xba, 0x65, 0x9f, 0x94, 0x72, 0x41, 0xbd, 0x40, 0x07, 0x68, 0x41, 0x5d,
	0x04, 0x5f, 0x29, 0x0f, 0x06, 0x96, 0xd1, 0x31, 0xdc, 0x96, 0x9f, 0x87, 0xce, 0x06, 0xb4, 0x9e,
	0x32, 0x35, 0xf4, 0x49, 0x09, 0xce, 0x14, 0x39, 0x2f, 0xb0, 0xfd, 0x9c, 0x85, 0x2a, 0x92, 0x89,
	0xd0, 0x09, 0x67, 0x3e, 0x7d, 0x64, 0xa4, 0xf4, 0x08, 0x40, 0x0d, 0x03, 0x39, 0xe8, 0xf7, 0xc6,
	0x00, 0x4d, 0x3f, 0x0f, 0xd1, 0x05, 0xf3, 0x2d, 0x49, 0x35, 0x49, 0x6b, 0xa5, 0x63, 0xb8, 0xeb,
	0xdd, 0xb6, 0x97, 0xb3, 0xf0, 0xee, 0xc6, 0x79, 0x7f, 0x5a, 0x77, 0x8e, 0xc0, 0x9c, 0x64, 0x70,
	0x07, 0xd6, 0x02, 0x21, 0x66, 0x58, 0x93, 0xc0, 0x39, 0x81, 0x46, 0x4e, 0x63, 0x01, 0xe1, 0x33,
	0xd8, 0xbb, 0xe1, 0x4c, 0x07, 0x09, 0x23, 0xf9, 0x40, 0x5a, 0x26, 0x91, 0xca, 0x49, 0x56, 0xc3,
	0x5e, 0x82, 0xf5, 0x77, 0xa0, 0x6a, 0xcd, 0x6a, 0x79, 0xcd, 0x29, 0x6c, 0xf9, 0x14, 0x11, 0xd3,
	0xf7, 0x3c, 0x5e, 0xb2, 0xc0, 0x03, 0x2c, 0xb7, 0x2e, 0x83, 0xee, 0x7e, 0x1b, 0x50, 0xef, 0x71,
	0x21, 0x52, 0x92, 0x78, 0x0d, 0xcd, 0xa9, 0xdc, 0x21, 0xe1, 0x61, 0x21, 0x5d, 0xc5, 0x0f, 0x6c,
	0x2c, 0xca, 0xb3, 0x77, 0xd5, 0xce, 0x0d, 0x7c, 0x85, 0xf6, 0xfc, 0x81, 0x78, 0x5c, 0xf4, 0xfe,
	0xa3, 0x96, 0xed, 0x2c, 0x6a, 0xc9, 0xe1, 0xb1, 0x0f, 0x50, 0x1c, 0x87, 0xfb, 0x65, 0x0a, 0x73,
	0xea, 0xd8, 0x07, 0xd5, 0xc5, 0x1c, 0xaa, 0xfb, 0x08, 0x9b, 0xd3, 0xb3, 0xfb, 0x2c, 0x26, 0xa5,
	0xb9, 0xc4, 0x2b, 0x30, 0x47, 0xee, 0x23, 0x89, 0xbb, 0xc5, 0x70, 0xd9, 0xb9, 0x76, 0x29, 0xff,
	0xcb, 0xa7, 0x35, 0xd7, 0x08, 0xcd, 0xb1, 0xed, 0x2f, 0x7e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x81,
	0xea, 0x3c, 0x01, 0x04, 0x03, 0x00, 0x00,
}
