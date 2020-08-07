package stringsvc

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
)

/*
使用 http&json作为传输协议组合
*/

var (
	svc    = stringService{}
	Logger = log.NewLogfmtLogger(os.Stderr)

	// 使用日志中间件记录传输层日志
	uppercase = loggingMiddleware(log.With(Logger, "method", "uppercase"))(makeUppercaseEndpoint(svc))
	count     = loggingMiddleware(log.With(Logger, "method", "count"))(makeCountEndpoint(svc))

	UppercaseHandler = httptransport.NewServer(
		uppercase,
		decodeUppercaseRequest,
		encodeResponse,
	)

	CountHandler = httptransport.NewServer(
		count,
		decodeCountRequest,
		encodeResponse,
	)
)

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
