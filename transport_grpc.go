package learn

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"golang.org/x/net/context"

	"github.com/briankassouf/kit/auth/jwt"
	"github.com/briankassouf/learn/pb"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func MakeGRPCServer(ctx context.Context, endpoints Endpoints, logger log.Logger) pb.UserServiceServer {
	options := []grpctransport.ServerOption{grpctransport.ServerErrorLogger(logger)}

	return &grpcServer{
		createUser: grpctransport.NewServer(
			ctx,
			endpoints.CreateUserEndpoint,
			DecodeGRPCCreateUserRequest,
			EncodeGRPCCreateUserResponse,
			append(options, grpctransport.ServerBefore(jwt.ToGRPCContext()))...,
		),
		getUser: grpctransport.NewServer(
			ctx,
			endpoints.GetUserEndpoint,
			DecodeGRPCGetUserRequest,
			EncodeGRPCGetUserResponse,
			options...,
		),
	}
}

type grpcServer struct {
	createUser grpctransport.Handler
	getUser    grpctransport.Handler
}

func (s *grpcServer) CreateUser(ctx context.Context, req *pb.CreateRequest) (*pb.UserResponse, error) {
	_, rep, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.UserResponse), nil
}

func (s *grpcServer) GetUser(ctx context.Context, req *pb.GetRequest) (*pb.UserResponse, error) {
	_, rep, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.UserResponse), nil
}

// DecodeGRPCCreateUserRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC create user request to a user-domain create user request. Primarily useful in a server.
func DecodeGRPCCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	return CreateUserRequest{
		User: &User{
			req.User.Id,
			req.User.FirstName,
			req.User.LastName,
			req.User.Email,
			req.User.Username,
		},
	}, nil
}

// DecodeGRPCGetUserRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC get user request to a user-domain get user request. Primarily useful in a server.
func DecodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetRequest)
	return GetUserRequest{Id: req.Id}, nil
}

// DecodeGRPCCreateUserResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC create user response to a user-domain create user response. Primarily useful in a client.
func DecodeGRPCCreateUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserResponse)
	return CreateUserResponse{
		User: &User{
			reply.User.Id,
			reply.User.FirstName,
			reply.User.LastName,
			reply.User.Email,
			reply.User.Username,
		},
		Err: nil,
	}, nil
}

// DecodeGRPCGetUserResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC get user response to a user-domain get user response. Primarily useful in a client.
func DecodeGRPCGetUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserResponse)
	return GetUserResponse{
		User: &User{
			reply.User.Id,
			reply.User.FirstName,
			reply.User.LastName,
			reply.User.Email,
			reply.User.Username,
		},
		Err: nil,
	}, nil
}

// EncodeGRPCCreateUserResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain Create User response to a gRPC Create User reply. Primarily useful in a server.
func EncodeGRPCCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CreateUserResponse)
	return &pb.UserResponse{
		User: &pb.User{
			resp.User.Id,
			resp.User.FirstName,
			resp.User.LastName,
			resp.User.Email,
			resp.User.Username,
		},
	}, nil
}

// EncodeGRPCGetUserResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain Get User response to a gRPC Get User reply. Primarily useful in a server.
func EncodeGRPCGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetUserResponse)
	return &pb.UserResponse{
		User: &pb.User{
			resp.User.Id,
			resp.User.FirstName,
			resp.User.LastName,
			resp.User.Email,
			resp.User.Username,
		},
	}, nil
}

// EncodeGRPCCreateUserRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Create User request to a gRPC Create User request. Primarily useful in a client.
func EncodeGRPCCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(CreateUserRequest)
	return &pb.CreateRequest{
		User: &pb.User{
			req.User.Id,
			req.User.FirstName,
			req.User.LastName,
			req.User.Email,
			req.User.Username,
		},
	}, nil
}

// EncodeGRPCGetUserRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Get User request to a gRPC Get User request. Primarily useful in a client.
func EncodeGRPCGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(GetUserRequest)
	return &pb.GetRequest{
		Id: req.Id,
	}, nil
}
