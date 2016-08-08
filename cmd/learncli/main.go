package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/briankassouf/learn"
	client "github.com/briankassouf/learn/client"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		grpcAddr = flag.String("grpc.addr", "", "gRPC (HTTP) address of addsvc")
		httpAddr = flag.String("http.addr", "", "http address")
		method   = flag.String("method", "create", "create, get")
	)
	flag.Parse()

	if len(flag.Args()) != 1 && *method == "get" {
		fmt.Fprintf(os.Stderr, "usage: learncli --method=get <id>\n")
		os.Exit(1)
	}

	if len(flag.Args()) != 5 && *method == "create" {
		fmt.Fprintf(os.Stderr, "usage: learncli --method=create <id> <first name> <last name> <email> <username>\n")
		os.Exit(1)
	}

	var service learn.UserService
	var err error
	if *httpAddr != "" {
		service, err = client.NewHTTP(*httpAddr, log.NewNopLogger())
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service = client.New(conn)

	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch *method {
	case "create":
		user := &learn.User{
			flag.Args()[0],
			flag.Args()[1],
			flag.Args()[2],
			flag.Args()[3],
			flag.Args()[4],
		}

		u, err := service.CreateUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(u)
	case "get":
		id := flag.Args()[0]

		u, err := service.GetUser(context.Background(), id)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(u)
	}
}
