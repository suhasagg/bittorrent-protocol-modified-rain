package subs

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/phonkee/go-pubsub"
	"github.com/pkg/errors"
	
)

type BittorrentSubscriptionSet struct {
        TorrentSession torrentsessionResetHub
	FileMetadata  filemetadataResetHub
	SystemProperty systempropertyResetHub
	IPAddressProperties ipaddressResetHub
	DataTransferProperties datatransferpropertiesresethub
}

func NewBittorrentSubscriptionSet() *BittorrentSubscriptionSet {
	s := &BittorrentSubscriptionSet{
		torrentsessionHub:      *newtorrentsessionResetHub(),
	        filemetadataHub:  *newfilemetadataResetHub(),
		systempropertyHub: *systempropertyResetHub(),
		ipaddresspropertiesHub: *ipaddressResetHub(),
		datatransferpropertiesHub: *datatransferpropertiesResetHub()
	}
	return s
}

func (s *BittorrentSubscriptionSet) AddSubscription(
	method string,
	filter bittorrent.bittorrentFilter,
	conn *websocket.Conn) (string, error) {
	var id string
	switch method {
	case TorrentSession:
		 id = s.torrentsessionhub.addSubscriber(filter, conn)
	case FileMetadata:
		 id = s.filemetadatahub.addSubscriber(filter, conn)
	case SystemProperty:
		 id = s.systempropertyhub.addSubscriber(filter, conn)
	//A peer has set of files, which are to be synchronised across certain number of peers, which have different set of files (replication only on peer consent), to ensure file availability is more uniform (very analogous CDN) 
	case IPAddressProperty:
	         id = s.ipaddresspropertyhub.addSubscriber(filter, conn)
        case DataTransferProperty:
	         id = s.datatransferpropertyhub.addSubscriber(filter, conn)
	default:
		return "", fmt.Errorf("unrecognised method %s", method)
	}
	return id, nil
}

func (s *BittorrentSubscriptionSet) EmitTorrentSessionEvent(data types.EventData) (err error) {
	return s.TorrentSessionHub.emitTorrentSessionEvent(data)
}

func (s *BittorrentSubscriptionSet) EmitFileMetadataEvent(data types.EventData) (err error) {
	return s.FileMetadataHub.emitFileMetaDataEvent(data)
}

func (s *BittorrentSubscriptionSet) EmitSystemPropertyEvent(data types.EventData) (err error) {
	return s.SystemPropertyHub.emitSystemPropertyEvent(data)
}

func (s *BittorrentSubscriptionSet) EmitIPAddressPropertyEvent(data types.EventData) (err error) {
	return s.IPAddressPropertyHub.emitIPAddressPropertyEvent(data)
}

func (s *BittorrentSubscriptionSet) EmitDataTransferPropertyEvent(data types.EventData) (err error) {
	return s.DataTransferPropertyHub.emitDataTransferPropertyEvent(data)
}

func (s *bittorrentSubscriptionSet) Remove(id string) {
	s.TorrentSessionHub.closeSubscription(id)
	s.FileMetadataHub.closeSubscription(id)
	s.SystemPropertyHub.closeSubscription(id)
	s.IPAddressPropertyHub.closeSubscription(id)
	s.DataTransferProperty.closeSubscription(id)

}

