package _go

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	do    func(context.Context) error
	clean func(err error)
	err   error
}

func (tk *Task) Valid() bool {
	return tk.clean != nil && tk.do != nil
}

// TaskGroup用于同时在后台启动一组同生命周期的多个任务（保证启动顺序与添加顺序一致），“同生共死”
type TaskGroup struct {
	tasks       []*Task
	tkBuf       *Task
	shareCtx    context.Context
	cancel      func()
	canceled    int32
	wg          sync.WaitGroup
	isScheduled bool
}

func NewTaskGroup() *TaskGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskGroup{
		shareCtx: ctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
	}
}

func (a *TaskGroup) Add(do func(context.Context) error) *TaskGroup {
	a.tkBuf = &Task{do: do}
	return a
}

func (a *TaskGroup) Interrupt(clean func(err error)) {
	if clean == nil {
		clean = func(err error) {}
	}
	a.tkBuf.clean = clean
	if a.tkBuf.Valid() {
		a.tasks = append(a.tasks, a.tkBuf)
		a.tkBuf = nil
	}
}

func (a *TaskGroup) schedule() {
	for _, f := range a.tasks {
		a.wg.Add(1)
		time.Sleep(time.Millisecond) // Guarantee schedule sequence
		go func(tk *Task) {
			defer func() {
				if err := recover(); err != nil {
					tk.err = fmt.Errorf("-------------panic: %v", err)
				}
				if tk.err != nil {
					a.cancelAll()
				}
				a.wg.Done() // call in last
			}()
			tk.err = tk.do(a.shareCtx)
		}(f)
	}
}

// Start start all the tasks as one goroutine per task
func (a *TaskGroup) Start() {
	if a.isScheduled {
		panic("go-util._go: all task have been scheduled!")
	}
	a.tkBuf = nil // clear buf
	a.schedule()
	a.isScheduled = true
}

func (a *TaskGroup) Wait() {
	a.wg.Wait()
}

// Run start all the tasks as one goroutine per task, then return until them done
func (a *TaskGroup) Run() {
	a.Start()
	a.wg.Wait()
}

func (a *TaskGroup) cancelAll() {
	if atomic.LoadInt32(&a.canceled) == 1 {
		return
	}
	atomic.CompareAndSwapInt32(&a.canceled, 0, 1)
	a.cancel()
	// reverse
	for i := len(a.tasks) - 1; i >= 0; i-- {
		tk := a.tasks[i]
		tk.clean(tk.err)
	}
}
