package client

import (
	"time"

	stdjwt "github.com/dgrijalva/jwt-go"
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	"github.com/briankassouf/kit/auth/jwt"
	"github.com/briankassouf/learn"
	"github.com/briankassouf/learn/pb"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// New returns an AddService backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn) learn.UserService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))
	options := []grpctransport.ClientOption{}
	jwtSigner := jwt.NewSigner("testSigningString1", stdjwt.SigningMethodHS256, stdjwt.MapClaims{})

	var createUserEndpoint endpoint.Endpoint
	{
		options = append(options, grpctransport.ClientBefore(jwt.FromGRPCContext()))

		createUserEndpoint = grpctransport.NewClient(
			conn,
			"UserService",
			"CreateUser",
			learn.EncodeGRPCCreateUserRequest,
			learn.DecodeGRPCCreateUserResponse,
			pb.UserResponse{},
			options...,
		).Endpoint()
		createUserEndpoint = limiter(createUserEndpoint)
		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateUser",
			Timeout: 30 * time.Second,
		}))(createUserEndpoint)
		createUserEndpoint = jwtSigner(createUserEndpoint)
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = grpctransport.NewClient(
			conn,
			"UserService",
			"GetUser",
			learn.EncodeGRPCGetUserRequest,
			learn.DecodeGRPCGetUserResponse,
			pb.UserResponse{},
			options...,
		).Endpoint()
		getUserEndpoint = limiter(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetUser",
			Timeout: 30 * time.Second,
		}))(getUserEndpoint)
	}

	return learn.Endpoints{
		CreateUserEndpoint: createUserEndpoint,
		GetUserEndpoint:    getUserEndpoint,
	}
}
