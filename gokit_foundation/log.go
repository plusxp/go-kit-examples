package gokit_foundation

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"time"
)

const SimpleLayout = "2006-01-02 15:04:05"

// 整个微服务中若要传递logger，只通过传递UnionLogger的方式
// 提醒：最好不要在项目中传递一种具体的logger struct，这会导致后续巨大的切换logger成本
type UnionLogger struct {
	// github.com/go-kit/kit/log 提供的所有logger都是Log(keyvals ...interface{})的调用形式，其实质是推崇这种kv形式的log
	// 这种形式的log丢失了一部分使用灵活度，但给日志分析带来了很大好处；
	Kvlgr log.Logger

	// 但没用过的开发者/团队或许会感到不太习惯；这没关系，在这里定义适用于你的logger接口，然后在项目中传递这个接口即可，
	// 比如logrus.FieldLogger，方便你仍可以调用log.Printf(...)
	Logruslgr logrus.FieldLogger
	// 继续添加你喜欢的logger
	Mylgr *CustomizeLogger
}

// 使用你习惯的logger初始化，在项目中调用对应logger的方法
func NewUnionLogger(kvLogger log.Logger, logrusLogger logrus.FieldLogger, mylgr *CustomizeLogger) *UnionLogger {
	return &UnionLogger{
		Kvlgr:     kvLogger,
		Logruslgr: logrusLogger,
		Mylgr:     mylgr,
	}
}

// 也可以自己实现log.Logger接口来自定义其底层行为
func NewKvLogger(logger log.Logger) log.Logger {
	var (
		ts                = log.TimestampFormat(time.Now, SimpleLayout)
		hommizationCaller = log.Valuer(func() interface{} {
			depth := 3
			_, file, line, _ := runtime.Caller(depth)
			s := fmt.Sprintf("%s:%d", file, line)
			ss := strings.Split(s, "/")
			// 限制路径深度为4
			// e.g. demo_project/gokit_foundation/gateway/gateway.go:38
			start := 0
			if len(ss)-3 > 0 {
				start = len(ss) - 3
			}
			return strings.Join(ss[start:], "/")
		})
	)

	var l log.Logger
	if logger != nil {
		l = logger
	} else {
		l = log.NewJSONLogger(os.Stderr)
	}

	l = log.With(l, "ts", ts)
	l = log.With(l, "caller", hommizationCaller)
	return l
}
