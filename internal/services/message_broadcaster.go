package services

import "sync"

type Message struct {
	Title string
	Body  string
}

type MessageBroadcaster struct {
	subscribers []chan Message
	mu          sync.Mutex
}

// An interface for Subscribers
type Subscriber interface {
	Run(wg *sync.WaitGroup)
}

func NewMessageBroadcaster() *MessageBroadcaster {
	return &MessageBroadcaster{
		subscribers: make([]chan Message, 0),
	}
}

func (b *MessageBroadcaster) Subscribe() <-chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Message, 10)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *MessageBroadcaster) Publish(msg Message) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subscribers {
		sub <- msg
	}
}

func (b *MessageBroadcaster) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subscribers {
		close(sub)
	}
}
