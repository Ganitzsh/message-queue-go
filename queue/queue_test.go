package queue_test

import (
	"testing"

	"github.com/ganitzsh/message-queue-go/queue"
	"github.com/ganitzsh/message-queue-go/types"
)

type Content struct{}

type TestMessage = types.Message[Content]

func newTestMessage() *TestMessage {
	return types.NewMessage("test-content", Content{})
}

func TestQueuePop(t *testing.T) {
	q := queue.NewQueue[Content]("test-queue", 200)
	defer q.Close()

	q.Send(newTestMessage())

	msg := <-q.Pop()
	if msg == nil {
		t.Fatalf("Expected message, got nil")
	}
}

func TestQueueWithOneWorker(t *testing.T) {
	q := queue.NewQueue[Content]("test-queue", 200)
	defer q.Close()

	q.Send(newTestMessage())

	worker := q.AddWorker(func(message *TestMessage) error {
		t.Logf("Received message: %v", message)
		return nil
	})

	msg := <-worker.Receive()
	if msg != nil {
		t.Fatalf("Expected no error, got error: %v", msg)
	}
}

func TestCloseAllWorkers(t *testing.T) {
	q := queue.NewQueue[Content]("test-queue", 200)

	q.AddWorker(func(message *TestMessage) error {
		t.Logf("Received message: %v", message)
		return nil
	})

	q.AddWorker(func(message *TestMessage) error {
		t.Logf("Received message: %v", message)
		return nil
	})

	q.Close()
}
