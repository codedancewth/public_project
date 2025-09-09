package utils

import (
	"fmt"
	"sync"
)

// 任务结构体
type Task struct {
	ID  int
	Job func() error
}

// 工作池结构体
type WorkerPool struct {
	NumWorkers int
	Tasks      chan Task
	WaitGroup  sync.WaitGroup
}

// 创建工作池
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		NumWorkers: numWorkers,
		Tasks:      make(chan Task),
	}
}

// 工作池启动函数
func (wp *WorkerPool) Start() {
	// 启动指定数量的goroutine，每个goroutine负责从任务通道中获取任务并执行
	for i := 0; i < wp.NumWorkers; i++ {
		go wp.worker()
	}
}

// 添加任务到工作池
func (wp *WorkerPool) AddTask(task Task) {
	wp.WaitGroup.Add(1)
	wp.Tasks <- task
}

// 等待所有任务完成
func (wp *WorkerPool) Wait() {
	wp.WaitGroup.Wait()
	close(wp.Tasks)
}

// 单个工作goroutine的函数
func (wp *WorkerPool) worker() {
	for task := range wp.Tasks {
		err := task.Job()
		if err != nil {
			fmt.Printf("Task %d failed: %v\n", task.ID, err)
		}
		wp.WaitGroup.Done()
	}
}

// 示例任务函数
func exampleTask() error {
	// 这里可以放置实际的任务逻辑
	fmt.Println("Executing task")
	return nil
}

//
//func main() {
//	// 创建一个包含3个工作goroutine的工作池
//	pool := NewWorkerPool(3)
//
//	// 启动工作池
//	pool.Start()
//
//	// 添加示例任务到工作池
//	for i := 0; i < 10; i++ {
//		task := Task{
//			ID:  i,
//			Job: exampleTask,
//		}
//		pool.AddTask(task)
//	}
//
//	// 等待所有任务完成
//	pool.Wait()
//
//	fmt.Println("All tasks completed")
//}
