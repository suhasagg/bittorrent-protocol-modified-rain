package announcer

import (
	"time"
	"github.com/cenkalti/rain/internal/logger"
)

// DHTAnnouncer runs a function periodically to announce the Torrent to DHT network.
type DHTAnnouncer struct {
	lastAnnounce   time.Time //When torrent was announced
	needMorePeers  bool //Need to grow peers
	needMorePeersC chan bool //Channel for co-ordination
	closeC         chan struct{}
	doneC          chan struct{}
}

// NewDHTAnnouncer returns a new DHTAnnouncer.
func NewDHTAnnouncer() *DHTAnnouncer {
	return &DHTAnnouncer{
		needMorePeersC: make(chan bool),
		closeC:         make(chan struct{}),
		doneC:          make(chan struct{}),
		needMorePeers:  true,
	}
}

// Close the announcer.
func (a *DHTAnnouncer) Close() {
	close(a.closeC)
	<-a.doneC
}

//As new peers are available they register in DHT
// NeedMorePeers signals the announcer goroutine to fetch more peers from DHT.
// Per announcement, one peer with most priority will be connected
func (a *DHTAnnouncer) NeedMorePeers(val bool) {
	select {
	case a.needMorePeersC <- val:  //If there is value available in channel - connect to peer with most priority in peer priority list
	case <-a.doneC:
	}
}

// Run the announcer. Invoke with go statement.
func (a *DHTAnnouncer) Run(announceFunc func(), interval, minInterval time.Duration, l logger.Logger) {
	defer close(a.doneC)

	timer := time.NewTimer(minInterval)
	defer timer.Stop()

	resetTimer := func() {
		if a.needMorePeers {
			timer.Reset(time.Until(a.lastAnnounce.Add(minInterval)))
		} else {
			timer.Reset(time.Until(a.lastAnnounce.Add(interval)))
		}
	}

	announce := func() {
		announceFunc()
		a.lastAnnounce = time.Now()
		resetTimer()
	}

	announce()
	for {
		select {
		//When time expires after d duration
		case <-timer.C:
		//announcement of more peers needed is done, records the announcement time
			announce() //Announcing after delta d
		case a.needMorePeers = <-a.needMorePeersC:
			resetTimer()
			return
		}
	}
}
