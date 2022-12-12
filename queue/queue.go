package queue

import (
	"fmt"
	"sync"

	"github.com/ganitzsh/message-queue-go/types"
	"github.com/ganitzsh/message-queue-go/worker"
)

type Queue[T interface{}] struct {
	sync.Mutex
	Name    string
	input   chan *types.Message[T]
	workers []*worker.Worker[T]
	wg      sync.WaitGroup
}

func NewQueue[T interface{}](name string, bufferSize uint32) *Queue[T] {
	return &Queue[T]{
		Name:  name,
		input: make(chan *types.Message[T], bufferSize),
	}
}

func (t *Queue[T]) Send(message *types.Message[T]) {
	t.input <- message
}

func (t *Queue[T]) AddWorker(handler func(message *types.Message[T]) error) *worker.Worker[T] {
	t.Lock()

	newWorker := worker.NewWorker(fmt.Sprintf("worker-%s-%d", t.Name, len(t.workers)+1), t.input)

	<-newWorker.Start(handler, &t.wg)

	t.workers = append(t.workers, newWorker)

	t.Unlock()

	return newWorker
}

func (t *Queue[T]) Pop() chan *types.Message[T] {
	return t.input
}

func (t *Queue[T]) Close() {
	close(t.input)
	t.wg.Wait()
}
