package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	jujuratelimit "github.com/juju/ratelimit"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/briankassouf/learn"
	"github.com/briankassouf/learn/pb"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/ratelimit"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	var (
		grpcAddr = flag.String("grpc.addr", ":8082", "gRPC (HTTP) listen address")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	// Metrics domain.

	// Metrics domain.
	var gets, creates metrics.Counter
	{
		// Business level metrics.
		creates = prometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: "learn",
			Name:      "user_create",
			Help:      "Total count of users created",
		}, []string{})
		gets = prometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: "learn",
			Name:      "user_get",
			Help:      "Total count of get operations",
		}, []string{})
	}
	var duration metrics.TimeHistogram
	{
		// Transport level metrics.
		duration = metrics.NewTimeHistogram(time.Nanosecond, prometheus.NewSummary(stdprometheus.SummaryOpts{
			Namespace: "learn",
			Name:      "request_duration_ns",
			Help:      "Request duration in nanoseconds.",
		}, []string{"method", "success"}))
	}

	// Business domain.
	var service learn.UserService
	{
		service = learn.NewBasicService()
		service = learn.ServiceLoggingMiddleware(logger)(service)
		service = learn.ServiceMetricsMiddleware(gets, creates)(service)
	}

	// Endpoint domain.
	var createUserEndpoint endpoint.Endpoint
	{
		createUserDuration := duration.With(metrics.Field{Key: "method", Value: "CreateUser"})
		createUserLogger := log.NewContext(logger).With("method", "CreateUser")
		limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(1, 1))

		createUserEndpoint = learn.MakeCreateUserEndpoint(service)
		createUserEndpoint = limiter(createUserEndpoint)
		createUserEndpoint = learn.EndpointLoggingMiddleware(createUserLogger)(createUserEndpoint)
		createUserEndpoint = learn.EndpointMetricsMiddleware(createUserDuration)(createUserEndpoint)
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserDuration := duration.With(metrics.Field{Key: "method", Value: "GetUser"})
		getUserLogger := log.NewContext(logger).With("method", "GetUser")
		limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(1, 1))

		getUserEndpoint = learn.MakeGetUserEndpoint(service)
		getUserEndpoint = limiter(getUserEndpoint)
		getUserEndpoint = learn.EndpointLoggingMiddleware(getUserLogger)(getUserEndpoint)
		getUserEndpoint = learn.EndpointMetricsMiddleware(getUserDuration)(getUserEndpoint)
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

	// Debug listener.
	go func() {
		logger := log.NewContext(logger).With("transport", "debug")

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", stdprometheus.Handler())

		logger.Log("addr", ":8080")
		errc <- http.ListenAndServe(":8080", m)
	}()

	// gRPC transport.
	go func() {
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := learn.MakeGRPCServer(ctx, endpoints, logger)
		s := grpc.NewServer()
		pb.RegisterUserServiceServer(s, srv)

		fmt.Println("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	fmt.Println("exit", <-errc)
}
