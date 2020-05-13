package pubsub

import (
	"context"
	"sync"

	"github.com/mistikel/api-websocket-go/utils"
)

type Manager interface {
	Subscribe(ctx context.Context, id string) <-chan string
	Unsubscribe(ctx context.Context, id string)
	Publish(ctx context.Context, message string)
	ResolveAllMessages(ctx context.Context) []string
}

func New() Manager {
	return &Pubsub{
		mutex:             &sync.Mutex{},
		messageChan:       make(chan string, 1),
		persistedMessages: []string{},
		subscriber:        make(map[string](chan string)),
	}
}

type Pubsub struct {
	mutex             *sync.Mutex
	messageChan       chan string
	subscriber        map[string]chan string
	persistedMessages []string
}

func (p *Pubsub) Unsubscribe(ctx context.Context, id string) {
	utils.InfoContext(ctx, "unsubscribe id: %s", id)
	p.mutex.Lock()
	close(p.subscriber[id])
	delete(p.subscriber, id)
	p.mutex.Unlock()
}

func (p *Pubsub) Subscribe(ctx context.Context, id string) <-chan string {
	utils.InfoContext(ctx, "try subscribe with id: %s", id)
	p.mutex.Lock()
	msgChan := make(chan string, 1)
	p.subscriber[id] = msgChan
	p.mutex.Unlock()

	return msgChan
}

func (p *Pubsub) Publish(ctx context.Context, message string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	utils.InfoContext(ctx, "try publish message %s", message)
	p.persistedMessages = append(p.persistedMessages, message)
	for _, c := range p.subscriber {
		c <- message
	}
}

func (p *Pubsub) ResolveAllMessages(ctx context.Context) []string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.persistedMessages
}
