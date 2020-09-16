package gokit_foundation

import (
	"github.com/go-kit/kit/log"
	"go-util/_util"
	"os"
	"strings"
	"time"
)

const SimpleLayout = "2006-01-02 15:04:05"

// 也可以自己实现log.Logger接口来自定义其底层行为
func NewLogger() log.Logger {
	var (
		logger            = log.NewJSONLogger(os.Stderr)
		ts                = log.TimestampFormat(time.Now, SimpleLayout)
		hommizationCaller = log.Valuer(func() interface{} {
			line := _util.FileWithLineNum(4)
			ss := strings.Split(line, "/")
			// 限制路径深度为3
			// e.g. gokit_foundation/gateway/gateway.go:38
			start := 0
			if len(ss)-3 > 0 {
				start = len(ss) - 3
			}
			return strings.Join(ss[start:], "/")
		})
	)

	logger = log.With(logger, "ts", ts)
	logger = log.With(logger, "caller", hommizationCaller)
	return logger
}
