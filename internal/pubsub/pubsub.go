package pubsub

import (
	"sync"

	"github.com/AdeshDeshmukh/crimson/internal/resp"
)

type Subscriber struct {
	channel chan resp.Value
}

func newSubscriber() *Subscriber {
	return &Subscriber{
		channel: make(chan resp.Value, 100),
	}
}

func (s *Subscriber) send(v resp.Value) {
	select {
	case s.channel <- v:
	default:
	}
}

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]*Subscriber
}

func New() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]*Subscriber),
	}
}

func (ps *PubSub) Subscribe(channels []string) (*Subscriber, []resp.Value) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := newSubscriber()

	for _, channel := range channels {
		ps.subscribers[channel] = append(ps.subscribers[channel], sub)
	}

	confirmations := make([]resp.Value, len(channels))
	for i, channel := range channels {
		confirmations[i] = resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{
				{Type: resp.BULK, Bulk: "subscribe"},
				{Type: resp.BULK, Bulk: channel},
				{Type: resp.INTEGER, Num: i + 1},
			},
		}
	}

	return sub, confirmations
}

func (ps *PubSub) Unsubscribe(sub *Subscriber, channels []string) []resp.Value {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, channel := range channels {
		subs := ps.subscribers[channel]
		for i, s := range subs {
			if s == sub {
				ps.subscribers[channel] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		if len(ps.subscribers[channel]) == 0 {
			delete(ps.subscribers, channel)
		}
	}

	confirmations := make([]resp.Value, len(channels))
	for i, channel := range channels {
		confirmations[i] = resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{
				{Type: resp.BULK, Bulk: "unsubscribe"},
				{Type: resp.BULK, Bulk: channel},
				{Type: resp.INTEGER, Num: 0},
			},
		}
	}

	return confirmations
}

func (ps *PubSub) Publish(channel, message string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs, exists := ps.subscribers[channel]
	if !exists {
		return 0
	}

	msg := resp.Value{
		Type: resp.ARRAY,
		Array: []resp.Value{
			{Type: resp.BULK, Bulk: "message"},
			{Type: resp.BULK, Bulk: channel},
			{Type: resp.BULK, Bulk: message},
		},
	}

	for _, sub := range subs {
		sub.send(msg)
	}

	return len(subs)
}

func (ps *PubSub) Receive(sub *Subscriber) resp.Value {
	return <-sub.channel
}