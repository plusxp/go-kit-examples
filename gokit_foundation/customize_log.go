package gokit_foundation

import (
	"fmt"
	"go-util/_util"
)

// 自定义logger，实现你需要的logger方法
type CustomizeLogger struct {
}

func (l *CustomizeLogger) tracePanicLine() string {
	line := _util.FileWithLineNum(3)
	return line
}

func (l *CustomizeLogger) PanicIfErr(err error, ignoreErrs []error, printText ...string) {
	if err != nil {
		var ig bool
		for _, e := range ignoreErrs {
			if e == err {
				ig = true
			}
		}

		if !ig {
			trace := l.tracePanicLine()
			panic(fmt.Sprintf("CustomizeLogger.PanicIfErr ERR:%s TRACE：%s", printText, trace))
		}
	}
}

// 在你希望panic时将panic代码行打印在第一行，以便快速找到出错位置
func (l *CustomizeLogger) Must(b bool, hint ...string) {
	if !b {
		trace := l.tracePanicLine()
		if len(hint) > 0 {
			panic(fmt.Sprintf("[from]:CustomizeLogger.Must [hint]:%s [trace]: %s", hint[0], trace))
		}
		panic(fmt.Sprintf("[from]:CustomizeLogger.Must [trace]: %s", trace))
	}
}
