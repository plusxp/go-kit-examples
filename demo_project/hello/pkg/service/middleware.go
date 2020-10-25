package service

import (
	"context"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(HelloService) HelloService

type loggingMiddleware struct {
	logger log.Logger
	next   HelloService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a HelloService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next HelloService) HelloService {
		return &loggingMiddleware{logger, next}
	}
}

func (l loggingMiddleware) SayHi(ctx context.Context, name string) (Response string, errCode pbcommon.R) {
	defer func() {
		l.logger.Log("method", "SayHi", "name", name, "Response", Response, "errCode", errCode)
	}()
	return l.next.SayHi(ctx, name)
}

func (l loggingMiddleware) MakeADate(c0 context.Context, p1 *pb.MakeADateRequest) (p0 *pb.MakeADateResponse, err error) {
	defer func() {
		l.logger.Log("method", "MakeADate", "p1", p1, "p0", p0)
	}()
	return l.next.MakeADate(c0, p1)
}

func (l loggingMiddleware) UpdateUserInfo(c0 context.Context, p1 *pb.UpdateUserInfoRequest) (p0 *pb.UpdateUserInfoResponse, e1 error) {
	defer func() {
		l.logger.Log("method", "UpdateUserInfo", "p1", p1, "p0", p0, "e1", e1)
	}()
	return l.next.UpdateUserInfo(c0, p1)
}
