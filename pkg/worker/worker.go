package worker

import (
	"golang.org/x/sync/errgroup"
)

// type worker[T any] struct {
// 	concurrency int
// 	tasks       []func() T
// }

// func NewWorker[T any](concurrency int) *worker[T] {
// 	if concurrency < 1 {
// 		concurrency = 1
// 	}

// 	return &worker[T]{
// 		concurrency: concurrency,
// 		tasks:       []func() T{},
// 	}
// }

// func (w *worker[T]) Add(fn func() T) {
// 	w.tasks = append(w.tasks, fn)
// }

// func (w *worker[T]) Run() []T {
// 	return ConcurrentRun(w.concurrency, w.tasks)
// }

// // 并发执行调度器
// func ConcurrentRun[T any](concurrency int, tasks []func() T) []T {
// 	var wg sync.WaitGroup

// 	semaphore := make(chan struct{}, concurrency)
// 	results := make([]T, len(tasks))

// 	for i, task := range tasks {
// 		wg.Add(1)
// 		semaphore <- struct{}{}

// 		go func(index int, taskFunc func() T) {
// 			defer func() {
// 				<-semaphore
// 				wg.Done()
// 			}()

// 			results[index] = taskFunc()
// 		}(i, task)
// 	}

// 	wg.Wait()

// 	return results
// }

// schedular

type worker struct {
	concurrency int
	tasks       []func() error
}

func NewWorker(concurrency int) *worker {
	if concurrency < 1 {
		concurrency = 1
	}

	return &worker{
		concurrency: concurrency,
		tasks:       []func() error{},
	}
}

func (w *worker) Add(fn func() error) {
	w.tasks = append(w.tasks, fn)
}

func (w *worker) Run() error {
	return concurrent(w.concurrency, w.tasks)
}

// 并发执行调度器
func concurrent(concurrency int, tasks []func() error) error {
	group := new(errgroup.Group)
	semaphore := make(chan struct{}, concurrency)

	for i, task := range tasks {
		semaphore <- struct{}{}

		_, task := i, task // capture loop variables
		group.Go(func() error {
			defer func() {
				<-semaphore
			}()

			if err := task(); err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}
