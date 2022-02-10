package subs

import (
	"github.com/phonkee/go-pubsub"
	"sync"
)

// hub implements Hub interface
// Remembers which subscribers messages have been published to
// and does not send repeat messages to any subscribers.
// Revert resets the memory of the subscribers that have received messages.
// Indexes information by id
type bittorrentResetHub struct {
	mutex   *sync.RWMutex
	unsent  map[string]bool
	clients map[string]pubsub.Subscriber
}

func neebittorrentResetHub() *ethResetHub {
	return &bittorrentResetHub{
		mutex:   &sync.RWMutex{},
		unsent:  make(map[string]bool),
		clients: make(map[string]pubsub.Subscriber),
	}
}

// CloseSubscriber removes subscriber from hub
func (h *bittorrentResetHub) CloseSubscriber(subscriber pubsub.Subscriber) {
	panic("should never be called")
}

func (h *bittorrentResetHub) closeSubscription(id string) {
	h.mutex.Lock()
	if sub, ok := h.clients[id]; ok {
		sub.Close()
	}
	delete(h.clients, id)
	delete(h.unsent, id)
	h.mutex.Unlock()
}

// Publish publishes message to subscribers
func (h *bittorrentResetHub) Publish(message pubsub.Message) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	count := 0
	for id, sub := range h.clients {
		if h.unsent[id] {
			if sub.Match(message.Topic()) {
				count += sub.Publish(message)
				h.unsent[id] = false
			}
		}
	}
	return count
}

// Subscribe adds subscription to topics and returns subscriber
func (h *bittorrentResetHub) Subscribe(_ ...string) pubsub.Subscriber {
	return nil
}

func (h *bittorrentResetHub) Reset() {
	h.mutex.Lock()
	for id := range h.unsent {
		h.unsent[id] = true
	}
	h.mutex.Unlock()
}
