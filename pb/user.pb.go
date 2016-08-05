// Code generated by protoc-gen-go.
// source: user.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	user.proto

It has these top-level messages:
	GetRequest
	CreateRequest
	UserResponse
	User
*/
package pb

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

type GetRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CreateRequest struct {
	User *User `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UserResponse struct {
	User *User `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *UserResponse) Reset()                    { *m = UserResponse{} }
func (m *UserResponse) String() string            { return proto.CompactTextString(m) }
func (*UserResponse) ProtoMessage()               {}
func (*UserResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UserResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type User struct {
	Id        string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	FirstName string `protobuf:"bytes,2,opt,name=firstName" json:"firstName,omitempty"`
	LastName  string `protobuf:"bytes,3,opt,name=lastName" json:"lastName,omitempty"`
	Email     string `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	Username  string `protobuf:"bytes,5,opt,name=username" json:"username,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*GetRequest)(nil), "pb.GetRequest")
	proto.RegisterType((*CreateRequest)(nil), "pb.CreateRequest")
	proto.RegisterType((*UserResponse)(nil), "pb.UserResponse")
	proto.RegisterType((*User)(nil), "pb.User")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for UserService service

type UserServiceClient interface {
	GetUser(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*UserResponse, error)
	CreateUser(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type userServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserServiceClient(cc *grpc.ClientConn) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := grpc.Invoke(ctx, "/pb.UserService/GetUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := grpc.Invoke(ctx, "/pb.UserService/CreateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceServer interface {
	GetUser(context.Context, *GetRequest) (*UserResponse, error)
	CreateUser(context.Context, *CreateRequest) (*UserResponse, error)
}

func RegisterUserServiceServer(s *grpc.Server, srv UserServiceServer) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.UserService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.UserService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _UserService_GetUser_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 237 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x92, 0xe1, 0xe2, 0x72, 0x4f,
	0x2d, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0xe2, 0xe3, 0x62, 0xca, 0x4c, 0x91, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb2, 0x94, 0x74, 0xb9, 0x78, 0x9d, 0x8b, 0x52, 0x13, 0x4b,
	0x52, 0x61, 0x0a, 0x64, 0xb8, 0x58, 0x40, 0x06, 0x80, 0x95, 0x70, 0x1b, 0x71, 0xe8, 0x15, 0x24,
	0xe9, 0x85, 0x02, 0xf9, 0x41, 0x60, 0x51, 0x25, 0x1d, 0x2e, 0x1e, 0x30, 0x2f, 0xb5, 0xb8, 0x20,
	0x3f, 0xaf, 0x38, 0x95, 0x80, 0xea, 0x26, 0x46, 0x2e, 0x16, 0x10, 0x17, 0xdd, 0x56, 0xa0, 0x36,
	0xce, 0xb4, 0xcc, 0xa2, 0xe2, 0x12, 0xbf, 0xc4, 0xdc, 0x54, 0x09, 0x26, 0xb0, 0x30, 0x42, 0x40,
	0x48, 0x8a, 0x8b, 0x23, 0x27, 0x11, 0x2a, 0xc9, 0x0c, 0x96, 0x84, 0xf3, 0x85, 0x44, 0xb8, 0x58,
	0x53, 0x73, 0x13, 0x33, 0x73, 0x24, 0x58, 0xc0, 0x12, 0x10, 0x0e, 0x48, 0x07, 0xc8, 0xc2, 0x3c,
	0x90, 0x0e, 0x56, 0x88, 0x0e, 0x18, 0xdf, 0xa8, 0x90, 0x8b, 0x1b, 0xe4, 0x86, 0xe0, 0xd4, 0xa2,
	0xb2, 0xcc, 0xe4, 0x54, 0x21, 0x5d, 0x2e, 0x76, 0x60, 0x70, 0x40, 0x5c, 0x05, 0x72, 0x2e, 0x22,
	0x6c, 0xa4, 0x04, 0xe0, 0xce, 0x87, 0x7a, 0x4f, 0x89, 0x41, 0xc8, 0x98, 0x8b, 0x0b, 0x12, 0x3e,
	0x60, 0x1d, 0x82, 0x20, 0x15, 0x28, 0xe1, 0x85, 0x4d, 0x53, 0x12, 0x1b, 0x38, 0xf4, 0x8d, 0x01,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xea, 0xd6, 0x95, 0x9e, 0x8b, 0x01, 0x00, 0x00,
}
