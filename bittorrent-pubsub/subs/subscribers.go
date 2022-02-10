package subs

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/bittorrent/plugin/types"
	"github.com/phonkee/go-pubsub"
)

type TorrentSessionDataSubscriber struct {
	wsSubscriber
	topic string
}

func newTorrentSessionDataSubscriber(hub pubsub.ResetHub, id, topic string, conn *websocket.Conn) pubsub.Subscriber {
	wsSub := newWsSubscriber(hub, conn, id)
	return TorrentSessionDataSubscriber{
		wsSubscriber: *wsSub,
		topic: topic,
	}
}

func (t TorrentSessionDataSubscriber) Match(topic string) bool {
	return t.topic == topic
}

type FileMetaDataSubscriber struct {
	wsSubscriber
	topic string
}

func newFileMetadataSubscriber(hub pubsub.ResetHub, id string, topic string, conn *websocket.Conn) FileMetadataSubscriber {
	wsSub := newWsSubscriber(hub, conn, id)
	return FileMetadataSubscriber{
		wsSubscriber: *wsSub,
		topic: topic,
	}
}

func (f FileMetadataSubscriber) Match(topic string) bool {
	  return f.topic == topic
}

type SystemPropertySubscriber struct {
	wsSubscriber
	topic string
}

func newSystemPropertySubscriber(hub pubsub.ResetHub, id string, topic string, conn *websocket.Conn) FileMetadataSubscriber {
	wsSub := newWsSubscriber(hub, conn, id)
	return SystemPropertySubscriber{
		wsSubscriber: *wsSub,
		topic: topic,
	}
}

func (f SystemPropertySubscriber) Match(topic string) bool {
          return f.topic == topic
}

type IPAddressPropertySubscriber struct {
	wsSubscriber
	topic string
}

func newIPAddressPropertySubscriber(hub pubsub.ResetHub, id string, topic string, conn *websocket.Conn) FileMetadataSubscriber {
	wsSub := newWsSubscriber(hub, conn, id)
	return IPAddressPropertySubscriber{
		wsSubscriber: *wsSub,
	        topic: topic,
	}
}

func (f IpAddressPropertySubscriber) Match(topic string) bool {
          return f.topic == topic
}


type DataTransferPropertySubscriber struct {
	wsSubscriber
	topic string
}

func newDataTransferPropertySubscriber(hub pubsub.ResetHub, id string, topic string, conn *websocket.Conn) FileMetadataSubscriber {
	wsSub := newWsSubscriber(hub, conn, id)
	return DataTransferPropertySubscriber{
		wsSubscriber: *wsSub,
		topic: topic,
	}
}

func (f DataTransferPropertySubscriber) Match(topic string) bool {
	 return f.topic == topic
}
