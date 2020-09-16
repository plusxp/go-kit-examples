package service

import (
	"context"
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

func (l loggingMiddleware) SayHi(ctx context.Context, name string) (reply string, errCode pbcommon.R) {
	defer func() {
		l.logger.Log("method", "SayHi", "name", name, "reply", reply, "errCode", errCode)
	}()
	return l.next.SayHi(ctx, name)
}

func (l loggingMiddleware) MakeADate(c0 context.Context, p1 *pb.MakeADateReq) (p0 *pb.MakeADateRsp) {
	defer func() {
		l.logger.Log("method", "MakeADate", "p1", p1, "p0", p0)
	}()
	return l.next.MakeADate(c0, p1)
}
