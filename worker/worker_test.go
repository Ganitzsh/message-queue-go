package worker_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/ganitzsh/message-queue-go/types"
	"github.com/ganitzsh/message-queue-go/worker"
)

type Content struct{}

type TestMessage = types.Message[Content]

func newTestMessage() *TestMessage {
	return types.NewMessage("test_content", Content{})
}

func TestWorkerHandlerNoError(t *testing.T) {
	worker := worker.NewWorker[Content]("solo-worker", nil)
	defer worker.Close()

	worker.Start(func(message *TestMessage) error {
		return nil
	}, nil)

	worker.Write(newTestMessage())

	err := <-worker.Receive()
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}
}

func TestWorkerHandlerError(t *testing.T) {
	worker := worker.NewWorker[Content]("solo-worker", nil)
	defer worker.Close()

	worker.Start(func(message *TestMessage) error {
		return errors.New("test error")
	}, nil)

	worker.Write(newTestMessage())

	err := <-worker.Receive()
	if err == nil {
		t.Fatalf("Expected error, got no error")
	}
}

func TestCannotStartMultipleRoutines(t *testing.T) {
	worker := worker.NewWorker[Content]("worker-1", nil)

	defer worker.Close()

	messageHandler := func(message *TestMessage) error {
		return nil
	}

	worker.Start(messageHandler, nil)
	worker.Start(messageHandler, nil)
}

func TestExternalInputTestExternalInput(t *testing.T) {
	externalInput := make(chan *TestMessage, 200)

	worker := worker.NewWorker("solo-worker", externalInput)
	defer worker.Close()

	worker.Start(func(message *TestMessage) error {
		return errors.New("test error")
	}, nil)

	worker.Write(newTestMessage())

	err := <-worker.Receive()
	if err == nil {
		t.Fatalf("Expected error, got no error")
	}
}

func TestReady(t *testing.T) {
	externalInput := make(chan *TestMessage, 200)

	worker := worker.NewWorker("solo-worker", externalInput)
	defer worker.Close()

	ready := <-worker.Start(func(message *TestMessage) error {
		return nil
	}, nil)

	if !ready {
		t.Fatalf("Expected ready, got not ready")
	}
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	worker := worker.NewWorker[Content]("solo-worker-waitgroup", nil)

	<-worker.Start(func(message *TestMessage) error {
		return nil
	}, &wg)

	worker.Close()

	wg.Wait()
}
