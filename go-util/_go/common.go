package _go

import (
	"context"
	"fmt"
	"sync"
)

type Setter interface {
	SetErr(err error)
	Err() error
}

type AsyncTask func(context.Context, Setter)

type setter struct {
	mutex sync.RWMutex
	err   error
}

func (t *setter) SetErr(err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.err == nil {
		t.err = err
	}
}

func (t *setter) Err() error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.err
}

// AsyncTask用于异步启动N个协程，并可以在当前协程与其同步
// ctx, ctxCancel用于在外部控制 【所有的task的生命周期】
type SafeAsyncTask struct {
	tasks         []AsyncTask
	setter        Setter
	ctx           context.Context
	cancelAllTask context.CancelFunc
	wg            sync.WaitGroup
	isScheduled   bool
	err           error // capture first err happened in these goroutines
}

func NewSafeAsyncTask(ctx context.Context, cancelFunc context.CancelFunc) *SafeAsyncTask {
	if ctx == nil {
		ctx, cancelFunc = context.WithCancel(context.Background())
	}
	return &SafeAsyncTask{
		ctx:           ctx,
		cancelAllTask: cancelFunc,
		wg:            sync.WaitGroup{},
		setter: &setter{
			mutex: sync.RWMutex{},
		},
	}
}

func (a *SafeAsyncTask) AddTask(f ...AsyncTask) {
	a.tasks = append(a.tasks, f...)
}

func (a *SafeAsyncTask) schedule() {
	for _, f := range a.tasks {
		a.wg.Add(1)
		go func(f func(ctx context.Context, tgr Setter)) {
			defer func() {
				if err := recover(); err != nil {
					a.setter.SetErr(fmt.Errorf("panic: %v", err))
				}
				// 当一个goroutine非正常结束，所有任务都应该退出
				// 如果不希望task结束后，即使发生错误也不要影响其他task，不要调用setter.SetErr()
				// 但是task发生panic一定会导致所有goroutine退出
				if a.setter.Err() != nil {
					a.cancelAllTask()
				}
				a.wg.Done()
			}()
			f(a.ctx, a.setter)
		}(f)
	}
}

func (a *SafeAsyncTask) Run() {
	if a.isScheduled {
		panic("go-util._go: all task have been scheduled!")
	}
	a.schedule()
	a.isScheduled = true
}

func (a *SafeAsyncTask) RunAndWait() {
	a.Run()
	a.wg.Wait()
}

func (a *SafeAsyncTask) Clear() {
	a.tasks = nil
}

func (a *SafeAsyncTask) Err() error {
	return a.setter.Err()
}
