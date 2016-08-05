package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/briankassouf/learn"
	"github.com/briankassouf/learn/pb"
	"github.com/go-kit/kit/endpoint"
)

func main() {
	var (
		grpcAddr = flag.String("grpc.addr", ":8082", "gRPC (HTTP) listen address")
	)
	flag.Parse()

	// Business domain.
	var service learn.UserService
	{
		service = learn.NewBasicService()
	}

	// Endpoint domain.
	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = learn.MakeCreateUserEndpoint(service)
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = learn.MakeGetUserEndpoint(service)
	}

	endpoints := learn.Endpoints{
		CreateUserEndpoint: createUserEndpoint,
		GetUserEndpoint:    getUserEndpoint,
	}

	// Mechanical domain.
	errc := make(chan error)
	ctx := context.Background()

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// gRPC transport.
	go func() {
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := learn.MakeGRPCServer(ctx, endpoints)
		s := grpc.NewServer()
		pb.RegisterUserServiceServer(s, srv)

		fmt.Println("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	fmt.Println("exit", <-errc)
}
