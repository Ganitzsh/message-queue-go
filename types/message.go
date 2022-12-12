package types

import "github.com/google/uuid"

type Message[T interface{}] struct {
	ackId string

	MessageType string
	Content     T
}

func NewMessage[T interface{}](messageType string, content T) *Message[T] {
	return &Message[T]{
		ackId:       uuid.New().String(),
		MessageType: messageType,
		Content:     content,
	}
}

func (m *Message[T]) Ack() error {
	return nil
}

func (m *Message[T]) Nack() error {
	return nil
}
