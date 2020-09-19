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

func (l *CustomizeLogger) PanicIfErr(err error, ignore ...error) {
	if err != nil {
		var ig bool
		for _, e := range ignore {
			if e == err {
				ig = true
			}
		}
		if !ig {
			trace := l.tracePanicLine()
			panic(fmt.Sprintf("CustomizeLogger.PanicIfErr TRACE：%s", trace))
		}
	}
}

// 在你希望panic时将panic行打印在第一行，以便快速找到出错位置
func (l *CustomizeLogger) Must(b bool) {
	if !b {
		trace := l.tracePanicLine()
		panic(fmt.Sprintf("CustomizeLogger.PanicIfErr TRACE：%s", trace))
	}
}
