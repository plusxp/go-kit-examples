package _go

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type AsyncTask struct {
	do    func(context.Context) error
	clean func(err error)
	err   error
}

func (tk *AsyncTask) Valid() bool {
	return tk.clean != nil && tk.do != nil
}

// AsyncTask用于异步启动N个协程，并可以在当前协程与其同步
// ctx, ctxCancel用于在外部控制 【所有的task的生命周期】
type SafeAsyncTask struct {
	tasks       []*AsyncTask
	tkBuf       *AsyncTask
	errs        chan error
	ctx         context.Context
	cancel      func()
	canceled    int32
	wg          sync.WaitGroup
	isScheduled bool
}

func NewSafeAsyncTask(ctx context.Context) *SafeAsyncTask {
	if ctx == nil {
		ctx, _ = context.WithCancel(context.Background())
	}
	return &SafeAsyncTask{
		ctx:  ctx,
		wg:   sync.WaitGroup{},
		errs: make(chan error),
	}
}

func (a *SafeAsyncTask) AddDo(do func(context.Context) error) *SafeAsyncTask {
	a.tkBuf = &AsyncTask{do: do}
	return a
}

func (a *SafeAsyncTask) AddClean(clean func(err error)) {
	if clean == nil {
		clean = func(err error) {}
	}
	a.tkBuf.clean = clean
	if a.tkBuf.Valid() {
		a.tasks = append(a.tasks, a.tkBuf)
		a.tkBuf = nil
	}
}

func (a *SafeAsyncTask) schedule() {
	for _, f := range a.tasks {
		a.wg.Add(1)
		time.Sleep(time.Millisecond) // Guarantee schedule sequence
		go func(tk *AsyncTask) {
			defer func() {
				if err := recover(); err != nil {
					tk.err = fmt.Errorf("panic: %v", err)
				}
				if tk.err != nil && atomic.LoadInt32(&a.canceled) == 0 {
					a.errs <- tk.err
					close(a.errs)
				}
				a.wg.Done() // call in last
			}()
			tk.err = tk.do(a.ctx)
		}(f)
	}
}

// Start start all the tasks as one goroutine per task
func (a *SafeAsyncTask) Start() {
	if a.isScheduled {
		panic("go-util._go: all task have been scheduled!")
	}
	a.tkBuf = nil // clear buf
	go a.onErr()
	a.schedule()
	a.isScheduled = true
}

func (a *SafeAsyncTask) Wait() {
	a.wg.Wait()
}

// Run start all the tasks as one goroutine per task, then return until them done
func (a *SafeAsyncTask) Run() {
	a.Start()
	a.wg.Wait()
}

func (a *SafeAsyncTask) onErr() {
	err := <-a.errs
	atomic.CompareAndSwapInt32(&a.canceled, 0, 1)
	a.cancel(err)
	// reverse
	for i := len(a.tasks) - 1; i >= 0; i-- {
		tk := a.tasks[i]
		tk.clean(tk.err)
	}
}
