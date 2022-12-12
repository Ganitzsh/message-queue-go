package worker

import (
	"log"
	"sync"

	"github.com/ganitzsh/message-queue-go/types"
)

const (
	DEFAULT_WORKER_OUTPUT_BUFFER = 100
	DEFAULT_WORKER_INPUT_BUFFER  = 100
)

type Worker[T interface{}] struct {
	sync.Mutex
	ID            string
	input         chan *types.Message[T]
	output        chan error
	ready         chan bool
	started       bool
	externalInput bool
}

func NewWorker[T interface{}](ID string, input chan *types.Message[T]) *Worker[T] {
	finalInput := input
	if input == nil {
		finalInput = make(chan *types.Message[T], DEFAULT_WORKER_INPUT_BUFFER)
	}

	return &Worker[T]{
		ID:            ID,
		input:         finalInput,
		output:        make(chan error, DEFAULT_WORKER_OUTPUT_BUFFER),
		ready:         make(chan bool, 1),
		started:       false,
		externalInput: input != nil,
	}
}

func (s *Worker[T]) Write(message *types.Message[T]) error {
	s.input <- message

	return nil
}

func (s *Worker[T]) Receive() chan error {
	return s.output
}

func (s *Worker[T]) Start(handler func(message *types.Message[T]) error, wg *sync.WaitGroup) chan bool {
	s.Lock()

	if s.started {
		log.Printf("[%s] Worker already started", s.ID)

		s.Unlock()

		return s.ready
	}

	s.started = true

	log.Printf("[%s] Worker started", s.ID)
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		if wg != nil {
			defer wg.Done()
		}

		s.ready <- s.started

		s.Unlock()

		for {
			message, more := <-s.input
			if !more {
				log.Printf("[%s] Worker stopped", s.ID)

				return
			}

			log.Printf("[%s] Received message: %v", s.ID, message.MessageType)

			if err := handler(message); err != nil {
				log.Printf("[%s] Error handling message: %v", s.ID, err)

				message.Nack()

				s.output <- err
			} else {
				message.Ack()

				s.output <- nil
			}
		}
	}()

	return s.ready
}

func (s *Worker[T]) Close() {
	if !s.externalInput {
		close(s.input)
	}

	close(s.output)
	close(s.ready)
}
