package broadcaster

import (
	"sync"
)

type Message struct {
	Title string
	Body  string
}

type Broadcaster struct {
	subscribers []chan Message
	mu          sync.Mutex
}

// An interface for Subscribers
type Subscriber interface {
	Run(wg *sync.WaitGroup)
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan Message, 0),
	}
}

func (b *Broadcaster) Subscribe() <-chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Message, 10)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster) Publish(msg Message) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subscribers {
		sub <- msg
	}
}

func (b *Broadcaster) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subscribers {
		close(sub)
	}
}
