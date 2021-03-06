// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/service.proto

/*
Package almanac is a generated protocol buffer package.

It is generated from these files:
	proto/service.proto
	proto/storage.proto

It has these top-level messages:
	AppendRequest
	AppendResponse
	IngestRequest
	IngestResponse
	SearchRequest
	SearchResponse
	LogEntry
	BleveIndex
	ChunkId
	Chunk
*/
package almanac

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

// A request to record append a log entry to an open chunk.
type AppendRequest struct {
	Entry *LogEntry `protobuf:"bytes,1,opt,name=entry" json:"entry,omitempty"`
}

func (m *AppendRequest) Reset()                    { *m = AppendRequest{} }
func (m *AppendRequest) String() string            { return proto.CompactTextString(m) }
func (*AppendRequest) ProtoMessage()               {}
func (*AppendRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AppendRequest) GetEntry() *LogEntry {
	if m != nil {
		return m.Entry
	}
	return nil
}

type AppendResponse struct {
}

func (m *AppendResponse) Reset()                    { *m = AppendResponse{} }
func (m *AppendResponse) String() string            { return proto.CompactTextString(m) }
func (*AppendResponse) ProtoMessage()               {}
func (*AppendResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// A request to ingest a single log entry into the system.
type IngestRequest struct {
	// A json object representing the entry to ingest.
	EntryJson string `protobuf:"bytes,1,opt,name=entry_json,json=entryJson" json:"entry_json,omitempty"`
}

func (m *IngestRequest) Reset()                    { *m = IngestRequest{} }
func (m *IngestRequest) String() string            { return proto.CompactTextString(m) }
func (*IngestRequest) ProtoMessage()               {}
func (*IngestRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *IngestRequest) GetEntryJson() string {
	if m != nil {
		return m.EntryJson
	}
	return ""
}

type IngestResponse struct {
}

func (m *IngestResponse) Reset()                    { *m = IngestResponse{} }
func (m *IngestResponse) String() string            { return proto.CompactTextString(m) }
func (*IngestResponse) ProtoMessage()               {}
func (*IngestResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// A request to search for log entries.
type SearchRequest struct {
	// A start time in epoch millisconds, inclusive. The search only returns
	// entries whose timestamp is at least this value.
	StartMs int64 `protobuf:"varint,2,opt,name=start_ms,json=startMs" json:"start_ms,omitempty"`
	// An end time in epoch milliseconds. If non-zero, only entries whose
	// timestamp is before this value are returned.
	EndMs int64 `protobuf:"varint,3,opt,name=end_ms,json=endMs" json:"end_ms,omitempty"`
	// A text-format query.
	Query string `protobuf:"bytes,4,opt,name=query" json:"query,omitempty"`
	// The maximum number of results to return.
	Num int32 `protobuf:"varint,5,opt,name=num" json:"num,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SearchRequest) GetStartMs() int64 {
	if m != nil {
		return m.StartMs
	}
	return 0
}

func (m *SearchRequest) GetEndMs() int64 {
	if m != nil {
		return m.EndMs
	}
	return 0
}

func (m *SearchRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *SearchRequest) GetNum() int32 {
	if m != nil {
		return m.Num
	}
	return 0
}

// The result of searching for log entries.
type SearchResponse struct {
	// All the entries which have matched the search.
	Entries []*LogEntry `protobuf:"bytes,2,rep,name=entries" json:"entries,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SearchResponse) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*AppendRequest)(nil), "almanac.AppendRequest")
	proto.RegisterType((*AppendResponse)(nil), "almanac.AppendResponse")
	proto.RegisterType((*IngestRequest)(nil), "almanac.IngestRequest")
	proto.RegisterType((*IngestResponse)(nil), "almanac.IngestResponse")
	proto.RegisterType((*SearchRequest)(nil), "almanac.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "almanac.SearchResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Appender service

type AppenderClient interface {
	// Appends an entry to an open chunk on this appender.
	Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error)
	// Executes a search on any open chunk(s) on this appender.
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type appenderClient struct {
	cc *grpc.ClientConn
}

func NewAppenderClient(cc *grpc.ClientConn) AppenderClient {
	return &appenderClient{cc}
}

func (c *appenderClient) Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error) {
	out := new(AppendResponse)
	err := grpc.Invoke(ctx, "/almanac.Appender/Append", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appenderClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := grpc.Invoke(ctx, "/almanac.Appender/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Appender service

type AppenderServer interface {
	// Appends an entry to an open chunk on this appender.
	Append(context.Context, *AppendRequest) (*AppendResponse, error)
	// Executes a search on any open chunk(s) on this appender.
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}

func RegisterAppenderServer(s *grpc.Server, srv AppenderServer) {
	s.RegisterService(&_Appender_serviceDesc, srv)
}

func _Appender_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppenderServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/almanac.Appender/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppenderServer).Append(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Appender_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppenderServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/almanac.Appender/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppenderServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Appender_serviceDesc = grpc.ServiceDesc{
	ServiceName: "almanac.Appender",
	HandlerType: (*AppenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Append",
			Handler:    _Appender_Append_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Appender_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}

// Client API for Ingester service

type IngesterClient interface {
	Ingest(ctx context.Context, in *IngestRequest, opts ...grpc.CallOption) (*IngestResponse, error)
}

type ingesterClient struct {
	cc *grpc.ClientConn
}

func NewIngesterClient(cc *grpc.ClientConn) IngesterClient {
	return &ingesterClient{cc}
}

func (c *ingesterClient) Ingest(ctx context.Context, in *IngestRequest, opts ...grpc.CallOption) (*IngestResponse, error) {
	out := new(IngestResponse)
	err := grpc.Invoke(ctx, "/almanac.Ingester/Ingest", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Ingester service

type IngesterServer interface {
	Ingest(context.Context, *IngestRequest) (*IngestResponse, error)
}

func RegisterIngesterServer(s *grpc.Server, srv IngesterServer) {
	s.RegisterService(&_Ingester_serviceDesc, srv)
}

func _Ingester_Ingest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngesterServer).Ingest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/almanac.Ingester/Ingest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngesterServer).Ingest(ctx, req.(*IngestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ingester_serviceDesc = grpc.ServiceDesc{
	ServiceName: "almanac.Ingester",
	HandlerType: (*IngesterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ingest",
			Handler:    _Ingester_Ingest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}

// Client API for Mixer service

type MixerClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type mixerClient struct {
	cc *grpc.ClientConn
}

func NewMixerClient(cc *grpc.ClientConn) MixerClient {
	return &mixerClient{cc}
}

func (c *mixerClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := grpc.Invoke(ctx, "/almanac.Mixer/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Mixer service

type MixerServer interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}

func RegisterMixerServer(s *grpc.Server, srv MixerServer) {
	s.RegisterService(&_Mixer_serviceDesc, srv)
}

func _Mixer_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MixerServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/almanac.Mixer/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MixerServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Mixer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "almanac.Mixer",
	HandlerType: (*MixerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _Mixer_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}

func init() { proto.RegisterFile("proto/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4d, 0x4b, 0xf3, 0x40,
	0x10, 0x26, 0x6f, 0xde, 0x4d, 0xdb, 0x91, 0x96, 0xba, 0x7e, 0xc5, 0x82, 0x50, 0x72, 0xb1, 0x20,
	0x44, 0x88, 0x17, 0x7b, 0xf0, 0xa0, 0xd0, 0x83, 0x62, 0x2f, 0xf1, 0x07, 0x94, 0xb5, 0x1d, 0x6a,
	0xc4, 0xee, 0xa6, 0xbb, 0x1b, 0xb1, 0x37, 0x7f, 0xba, 0xec, 0x6e, 0x56, 0x12, 0xf4, 0xe4, 0x6d,
	0x66, 0x9f, 0x79, 0x3e, 0x66, 0x12, 0x38, 0x28, 0xa5, 0xd0, 0xe2, 0x52, 0xa1, 0x7c, 0x2f, 0x96,
	0x98, 0xda, 0x8e, 0x76, 0xd8, 0xdb, 0x86, 0x71, 0xb6, 0x1c, 0x79, 0x54, 0x0b, 0xc9, 0xd6, 0x35,
	0x9a, 0x5c, 0x43, 0xff, 0xb6, 0x2c, 0x91, 0xaf, 0x72, 0xdc, 0x56, 0xa8, 0x34, 0x3d, 0x07, 0x82,
	0x5c, 0xcb, 0x5d, 0x1c, 0x8c, 0x83, 0xc9, 0x5e, 0xb6, 0x9f, 0xd6, 0xf4, 0xf4, 0x51, 0xac, 0x67,
	0x06, 0xc8, 0x1d, 0x9e, 0x0c, 0x61, 0xe0, 0x99, 0xaa, 0x14, 0x5c, 0x61, 0x92, 0x42, 0xff, 0x9e,
	0xaf, 0x51, 0x69, 0xaf, 0x75, 0x06, 0x60, 0x67, 0x17, 0xaf, 0x4a, 0x70, 0x2b, 0xd8, 0xcb, 0x7b,
	0xf6, 0xe5, 0x41, 0x09, 0x6e, 0x14, 0xfc, 0x7c, 0xad, 0x50, 0x40, 0xff, 0x09, 0x99, 0x5c, 0xbe,
	0x78, 0x85, 0x53, 0xe8, 0x2a, 0xcd, 0xa4, 0x5e, 0x6c, 0x54, 0xfc, 0x6f, 0x1c, 0x4c, 0xc2, 0xbc,
	0x63, 0xfb, 0xb9, 0xa2, 0x47, 0x10, 0x21, 0x5f, 0x19, 0x20, 0xb4, 0x00, 0x41, 0xbe, 0x9a, 0x2b,
	0x7a, 0x08, 0x64, 0x5b, 0xa1, 0xdc, 0xc5, 0xff, 0xad, 0x9d, 0x6b, 0xe8, 0x10, 0x42, 0x5e, 0x6d,
	0x62, 0x32, 0x0e, 0x26, 0x24, 0x37, 0x65, 0x72, 0x03, 0x03, 0x6f, 0xe5, 0xcc, 0xe9, 0x05, 0x74,
	0x4c, 0xb6, 0x02, 0x8d, 0x55, 0xf8, 0xfb, 0xee, 0x7e, 0x22, 0xfb, 0x0c, 0xa0, 0xeb, 0xd6, 0x47,
	0x49, 0xa7, 0x10, 0xb9, 0x9a, 0x1e, 0x7f, 0x53, 0x5a, 0x57, 0x1d, 0x9d, 0xfc, 0x78, 0xaf, 0x4d,
	0xa7, 0x10, 0xb9, 0x18, 0x0d, 0x6a, 0xeb, 0x04, 0x0d, 0x6a, 0x3b, 0x6f, 0x36, 0x83, 0xae, 0x3b,
	0x9f, 0x4b, 0xe0, 0xea, 0x86, 0x4c, 0xeb, 0x5b, 0x34, 0x64, 0xda, 0x37, 0xcf, 0xee, 0x80, 0xcc,
	0x8b, 0x0f, 0xa7, 0xf1, 0xc7, 0x28, 0xcf, 0x91, 0xfd, 0x99, 0xae, 0xbe, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xb9, 0xab, 0xa2, 0x24, 0x81, 0x02, 0x00, 0x00,
}
