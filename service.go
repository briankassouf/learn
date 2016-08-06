package learn

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
)

type UserService interface {
	CreateUser(cxt context.Context, user *User) (*User, error)
	GetUser(cxt context.Context, id string) (*User, error)
}

type basicService struct {
	Users map[string]*User
}

func NewBasicService() UserService {
	return basicService{
		make(map[string]*User),
	}
}

func (s basicService) CreateUser(_ context.Context, user *User) (*User, error) {
	s.Users[user.Id] = user

	return user, nil
}

func (s basicService) GetUser(_ context.Context, id string) (*User, error) {
	user, ok := s.Users[id]
	if !ok {
		return nil, errors.New("Could not find user")
	}

	return user, nil
}

type Middleware func(UserService) UserService

func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next UserService) UserService {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   UserService
}

func (mw serviceLoggingMiddleware) CreateUser(ctx context.Context, u *User) (user *User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CreateUser",
			"user", fmt.Sprintf("%v", u), "result", fmt.Sprintf("%v", user), "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.CreateUser(ctx, u)
}

func (mw serviceLoggingMiddleware) GetUser(ctx context.Context, id string) (user *User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUser",
			"id", id, "result", fmt.Sprintf("%v", user), "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetUser(ctx, id)
}

func ServiceMetricsMiddleware(gets metrics.Counter, creates metrics.Counter) Middleware {
	return func(next UserService) UserService {
		return serviceMetricsMiddleware{
			gets:    gets,
			creates: creates,
			next:    next,
		}
	}
}

type serviceMetricsMiddleware struct {
	gets    metrics.Counter
	creates metrics.Counter
	next    UserService
}

func (mw serviceMetricsMiddleware) CreateUser(ctx context.Context, u *User) (*User, error) {
	defer mw.creates.Add(1)
	return mw.next.CreateUser(ctx, u)
}

func (mw serviceMetricsMiddleware) GetUser(ctx context.Context, id string) (*User, error) {
	defer mw.gets.Add(1)
	return mw.next.GetUser(ctx, id)
}

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Username  string
}
