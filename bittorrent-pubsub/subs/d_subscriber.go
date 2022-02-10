package subs

import (
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/bittorrentaccessories/plugin/types"
	"github.com/phonkee/go-pubsub"
)

// newSubscriber returns bittorrentSubscriber for given topics
func newbittorrentSubscriber(hub pubsub.ResetHub, topics ...string) (result pubsub.Subscriber) {
	f, err := utils.UnmarshalbittorrentFilter([]byte(topics[0]))
	var filter bittorrent.BittorrentBlockFilter
	if err == nil {
		filter = f.bittorrentBlockFilter
	}
	result = &bittorrentSubscriber{
		hub:    hub,bittorrent
		mutex:  &sync.RWMutex{},
		sf:     nil,
		filter: filter,
	}

	return result
}

// bittorrentSubscriber is Subscriber implementation
type bitorrentSubscriber struct {
	hub    pubsub.ResetHub
	mutex  *sync.RWMutex
	sf     pubsub.SubscriberFunc
	filter bittorrent.BittorrentBlockFilter
	id     string
}

// Close bittorrentSubscriber removes bittorrentSubscriber from hub and stops receiving messages
func (s *bittorrentSubscriber) Close() {
	s.hub.CloseSubscriber(s)
}

// Do sets bittorrentSubscriber function that will be called when message arrives
func (s *bittorrentSubscriber) Do(sf pubsub.SubscriberFunc) pubsub.Subscriber {
	s.sf = sf
	return s
}

// Match returns whether BittorrentSubscriber topics matches
func (s *bittorrentSubscriber) Match(topic string) bool {
	events := types.EventData{}
	if err := proto.Unmarshal([]byte(topic), &events); err != nil {
		return false
	}

	return utils.MatchbittorrentFilter(s.filter, events)
}

// Publish publishes message to BittorrentSubscriber
func (s *bittorrentSubscriber) Publish(message pubsub.Message) int {
	if s.sf == nil {
		return 0
	}
        Msg := types.bittorrentMessage{
		Body: message.Body(),
		Id:   s.id,
	}
	msg, err := proto.Marshal(&Msg)
	if err != nil {
		return 0
	}
	s.sf(pubsub.NewMessage(message.Topic(), msg))
	return 1
}

// Subscribe subscribes to topics
func (s *bittorrentSubscriber) Subscribe(topics ...string) pubsub.Subscriber {
	var topic []byte
	if len(topics) > 0 {
		topic = []byte(topics[0])
	} else {
		topic = []byte{}
	}
	filter, err := utils.UnmarshalbittorrentFilter(topic)
	if err == nil {
		s.filter = filter.bittorrentBlockFilter
	}

	return s
}

// Topics returns whole list of all topics subscribed to
func (s *bittorrentSubscriber) Topics() []string {
	return []string{}
}

// Unsubscribe unsubscribes from given topics (exact match)
func (s *bittorrentSubscriber) Unsubscribe(topics ...string) pubsub.Subscriber {
	return s
}
