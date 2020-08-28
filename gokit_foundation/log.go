package gokit_foundation

import (
	"github.com/go-kit/kit/log"
	"os"
)

// 也可以自己实现log.Logger接口来自定义其底层行为
func NewLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}
