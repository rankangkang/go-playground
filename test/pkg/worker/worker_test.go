package worker_test

import (
	"fmt"
	"testing"
	"time"

	"playground/pkg/worker"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	len := 20

	worker := worker.NewWorker(4)
	result := []int{}
	for i := 0; i < len; i++ {
		var ti = i
		fn := func() error {
			time.Sleep(10 * time.Millisecond * 100)
			result = append(result, ti)
			return nil
		}
		worker.Add(fn)
	}

	err := worker.Run()
	assert.Nil(t, err)
	fmt.Println(result)
}

func TestNoWorker(t *testing.T) {
	len := 20
	result := []int{}
	for i := 0; i < len; i++ {
		var ti = i
		fn := func() (int, error) {
			time.Sleep(10 * time.Millisecond * 100)
			return ti, nil
		}

		ret, err := fn()
		if err != nil {
			assert.Nil(t, err)
		}

		result = append(result, ret)
	}

	fmt.Println(result)
}
