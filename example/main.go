package main

import (
	"log"

	"github.com/ganitzsh/message-queue-go/queue"
	"github.com/ganitzsh/message-queue-go/types"
)

type MessageContent struct {
	Foo string
}

type MyMessage = types.Message[MessageContent]

func newMyMessage() *MyMessage {
	return types.NewMessage("my-unique-message-type", MessageContent{
		Foo: "bar",
	})
}

func messageHandler(message *MyMessage) error {
	log.Printf("Received a fancy message saying: %s", message.Content.Foo)

	return nil
}

func main() {
	q := queue.NewQueue[MessageContent]("some-fancy-name", 200)
	defer q.Close()

	q.AddWorker(messageHandler)
	q.AddWorker(messageHandler)
	q.AddWorker(messageHandler)
	q.AddWorker(messageHandler)

	msg := newMyMessage()

	for i := 0; i < 100; i++ {
		q.Send(msg)
	}
}
