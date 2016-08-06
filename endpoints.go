package learn

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Endpoints struct {
	CreateUserEndpoint endpoint.Endpoint
	GetUserEndpoint    endpoint.Endpoint
}

// CreateUser implements Service. Primarily useful in a client.
func (e Endpoints) CreateUser(ctx context.Context, user *User) (*User, error) {
	request := CreateUserRequest{User: user}
	response, err := e.CreateUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(CreateUserResponse).User, nil
}

// GetUser implements Service. Primarily useful in a client.
func (e Endpoints) GetUser(ctx context.Context, id string) (*User, error) {
	request := GetUserRequest{Id: id}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(GetUserResponse).User, nil
}

func MakeCreateUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userRequest := request.(CreateUserRequest)
		user, err := s.CreateUser(ctx, userRequest.User)

		return CreateUserResponse{
			User: user,
			Err:  err,
		}, nil
	}
}

func MakeGetUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userRequest := request.(GetUserRequest)
		user, err := s.GetUser(ctx, userRequest.Id)

		return GetUserResponse{
			User: user,
			Err:  err,
		}, nil
	}
}

func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func EndpointMetricsMiddleware(duration metrics.TimeHistogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				f := metrics.Field{Key: "success", Value: fmt.Sprint(err == nil)}
				duration.With(f).Observe(time.Since(begin))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

type CreateUserRequest struct {
	User *User
}

type CreateUserResponse struct {
	User *User
	Err  error
}

type GetUserRequest struct {
	Id string
}

type GetUserResponse struct {
	User *User
	Err  error
}
