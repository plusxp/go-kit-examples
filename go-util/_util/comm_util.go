package _util

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"go-util/_go"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func ListenSignalTask(ctx context.Context, cancel context.CancelFunc, logger log.Logger, onClose func()) _go.AsyncTask {
	return func(context.Context, _go.Setter) {
		logger.Log("NewSafeAsyncTask", "ListenSignal")
		sc := make(chan os.Signal)
		signal.Notify(sc,
			syscall.SIGINT,  // 键盘中断
			syscall.SIGTERM, // 软件终止
		)
		select {
		case s := <-sc:
			fmt.Fprint(os.Stdout, "\n")
			logger.Log("ListenSignalTask", "===================== Closing ======================")
			logger.Log("recv-signal", s, "Task", "done!")
		case <-ctx.Done():
			fmt.Fprint(os.Stdout, "\n")
			logger.Log("ListenSignalTask", "===================== Closing ======================")
			logger.Log("ctx.Done", "OK", "Task", "done!")
		}
		onClose()
		signal.Stop(sc)
		cancel()
	}
}

func InCollection(elem interface{}, coll []interface{}) bool {
	for _, e := range coll {
		if e == elem {
			return true
		}
	}
	return false
}

func PanicIfErr(err interface{}, ignoreErrs []error, printText ...string) {
	if err != nil {
		for _, e := range ignoreErrs {
			if err == e {
				return
			}
		}
		if len(printText) > 0 {
			panic(printText[0])
		}
		panic(err)
	}
}

func AnyErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

func Must(condition bool, err error) {
	if !condition {
		panic(err)
	}
}

func If(condition bool, then func(), _else ...func()) {
	if condition {
		if then != nil {
			then()
		}
	} else {
		for _, f := range _else {
			f()
		}
	}
}

type SvcWithClose interface {
	Close() error
}

func CloseSvcSafely(manySvc []SvcWithClose) []error {
	var (
		errs []error
		err  error
	)
	for _, s := range manySvc {
		if reflect.ValueOf(s).IsNil() {
			continue
		}
		if err = s.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

var shanghai, _ = time.LoadLocation("Asia/Shanghai")
var simpleLayout = "2006-01-02 15:04:05"

func LoadShanghaiTimeFromStr(s string) (time.Time, error) {
	return time.ParseInLocation(simpleLayout, s, shanghai)
}

// 获取指定函数的名称, split:分割符，`.`获取纯函数名， `/`获取带pkg的函数名，如 _util.GetFuncName
func GetFuncName(funcObj interface{}, split ...string) string {
	fn := runtime.FuncForPC(reflect.ValueOf(funcObj).Pointer()).Name()
	if len(split) > 0 {
		fs := strings.Split(fn, split[0])
		return fs[len(fs)-1]
	}
	return fn
}

// 当前运行的函数名, split:分割符，不传就是获取全路径的函数名称
// split 传入 `.`获取纯函数名， `/`获取带pkg的函数名，如 _util.GetRunningFuncName
func GetRunningFuncName(split ...string) string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	fn := runtime.FuncForPC(pc[0]).Name()

	if len(split) > 0 {
		fs := strings.Split(fn, split[0])
		return fs[len(fs)-1]
	}
	return fn
}

// skip=1 为调用者位置，skip=2为调用者往上一层的位置，以此类推
// return-example: /develop/go/test_go/tmp_test.go:88
func FileWithLineNum(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%v:%v", file, line)
}
