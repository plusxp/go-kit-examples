package gokit_foundation

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"os"
	"runtime"
	"strings"
	"time"
)

const TimeCommonLayout = "2006-01-02 15:04:05"

// 项目中使用一种logger即可，不建议同时使用多个不同风格的logger
type Logger struct {
	log.Logger
	//logrus.FieldLogger
	*CustomizeLogger // 自定义logger，作为扩展
}

func NewLogger(logger log.Logger) *Logger {
	return &Logger{
		Logger:          NewKvLogger(logger),
		CustomizeLogger: new(CustomizeLogger),
	}
}

// 实现log.Logger接口来自定义其底层行为
func NewKvLogger(logger log.Logger) log.Logger {
	if logger != nil {
		return logger
	}
	var (
		ts                = log.TimestampFormat(time.Now, TimeCommonLayout)
		hommizationCaller = log.Valuer(func() interface{} {
			/*
			 获得更简洁的caller位置
			*/
			depth := 3 // 如果你不知道其含义，不要修改这个值
			_, file, line, _ := runtime.Caller(depth)
			s := fmt.Sprintf("%s:%d", file, line)
			ss := strings.Split(s, "/")
			// 限制路径层数为3，如果需要更完整的路径，增加即可
			// e.g. gokit_foundation/gateway/gateway.go:38
			layer := 3
			start := 0
			if len(ss)-layer > 0 {
				start = len(ss) - layer
			}
			return strings.Join(ss[start:], "/")
		})
	)

	var l log.Logger
	l = log.NewLogfmtLogger(os.Stdout)
	l = log.With(l, "ts", ts)
	l = log.With(l, "caller", hommizationCaller)
	return l
}
