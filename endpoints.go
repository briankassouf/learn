package learn

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUserEndpoint endpoint.Endpoint
	GetUserEndpoint    endpoint.Endpoint
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
