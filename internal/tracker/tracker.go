// Package tracker provides support for announcing torrents to HTTP and UDP trackers.
package tracker

import (
	"context"
	"errors"
	"net"
	"time"
)

/*
udp://tracker.coppersurfer.tk:6969/announce
udp://tracker.opentrackr.org:1337/announce
wss://tracker.openwebtorrent.com

secure web socket based tracker module + publish / subscribe hub support based on web sockets (peer, peer subgroup data exchange)

*/

// Tracker tracks the IP address of peers of a Torrent swarm.
// Torrent tracker (If a new peer is available which has file, it will do announcement to the tracker, (peer list is growing for that torrent)
type Tracker interface {
	// Announce transfer to the tracker.
	// Announce should be called periodically with the interval returned in AnnounceResponse.
	// Announce should also be called on specific events.
	Announce(ctx context.Context, req AnnounceRequest) (*AnnounceResponse, error)

	// URL of the tracker.
	URL() string
}

// AnnounceRequest contains the parameters that are sent in an announce request to trackers.
type AnnounceRequest struct {
	Torrent Torrent
	Event   Event
	NumWant int
}

// AnnounceResponse contains fields from a response to announce request.
type AnnounceResponse struct {
	Interval       time.Duration
	MinInterval    time.Duration
	Leechers       int32
	Seeders        int32
	WarningMessage string
	Peers          []*net.TCPAddr
}

// ErrDecode is returned from Tracker.Announce method when there is problem with the encoding of response.
var ErrDecode = errors.New("cannot decode response")

// Error is the string that is sent by the tracker from announce or scrape.
type Error struct {
	FailureReason string
	RetryIn       time.Duration
}

func (e *Error) Error() string { return e.FailureReason }
