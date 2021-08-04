package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Job struct {
	data int
}

type WorkerPool struct {
	jobs    chan Job
	workers int
	Done    chan struct{}
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Job, 50),
		workers: workers,
		Done:    make(chan struct{}),
	}
}

func (w *WorkerPool) Run(stop <-chan struct{}) {
	wg := sync.WaitGroup{}
	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go w.work(&wg)
	}
	log.Print("all worker started")
	<-stop
	// 接收到停止信号，先停止接受新任务
	log.Println("stop job chan")
	close(w.jobs)
	close(w.Done)
	// 等待所有 worker 处理完手头任务并安全退出
	wg.Wait()
	log.Println("all worker stopped")
}

func (w *WorkerPool) Add(j Job) {
	w.jobs <- j
}

func (w *WorkerPool) Stop() {}

func (w *WorkerPool) work(wg *sync.WaitGroup) {
	defer wg.Done()
	// 使用 for range 读取 chan 数据，chan 关闭时会自动退出循环
	for job := range w.jobs {
		w.process(&job)
	}
	log.Println("job chan closed, stop worker")
}

func (w *WorkerPool) process(j *Job) {
	// TODO: 单个任务执行的超时时间
	time.Sleep(500 * time.Millisecond)
	log.Printf("processed job [%d], result [%d]", j.data, j.data*j.data)
}

func main() {
	stop := make(chan struct{})
	w := NewWorkerPool(10)
	go func() {
		for i := 0; i < 50; i++ {
			w.Add(Job{data: i})
		}
	}()

	// 等待信号
	c := make(chan os.Signal, 2)
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
	signal.Notify(c, shutdownSignals...)
	go func() {
		log.Println("catch signal", <-c)
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	w.Run(stop)
}
