package addrlist

import (
	"net"
	"sort"
	"time"

	"github.com/cenkalti/rain/internal/blocklist"
	"github.com/cenkalti/rain/internal/externalip"
	"github.com/cenkalti/rain/internal/peerpriority"
	"github.com/cenkalti/rain/internal/peersource"
	"github.com/google/btree"
)

// AddrList contains peer addresses that are ready to be connected.
type AddrList struct {
	peerByTime     []*peerAddr
	//Application of btree APIs
	peerByPriority *btree.BTree

	maxItems   int
	listenPort int
	clientIP   *net.IP
	//If peer is there in block list, it wont be part of p2p network
	blocklist  *blocklist.Blocklist

	countBySource map[peersource.Source]int
}

// New returns a new AddrList.
func New(maxItems int, blocklist *blocklist.Blocklist, listenPort int, clientIP *net.IP) *AddrList {
	return &AddrList{
		peerByPriority: btree.New(2),

		maxItems:      maxItems,
		listenPort:    listenPort,
		clientIP:      clientIP,
		blocklist:     blocklist,
		countBySource: make(map[peersource.Source]int),
	}
}

// Reset empties the address list.
func (d *AddrList) Reset() {
	d.peerByTime = nil
	d.peerByPriority.Clear(false)
	d.countBySource = make(map[peersource.Source]int)
}

// Len returns the number of addresses in the list.
func (d *AddrList) Len() int {
	return d.peerByPriority.Len()
}

// LenSource returns the number of addresses for a single source.
func (d *AddrList) LenSource(s peersource.Source) int {
	return d.countBySource[s]
}

// Pop returns the next address. The returned address is removed from the list.
func (d *AddrList) Pop() (*net.TCPAddr, peersource.Source) {
	//Delete the node farthest from the peer calculated according to peer priority algorithm
	item := d.peerByPriority.DeleteMax()
	if item == nil {
		return nil, 0
	}
	p := item.(*peerAddr)
	d.peerByTime[p.index] = nil
	d.countBySource[p.source]--
	return p.addr, p.source
}
//Connection established - peer removed, freshly added

// Push adds a new address to the list. Does nothing if the address is already in the list.
func (d *AddrList) Push(addrs []*net.TCPAddr, source peersource.Source) {
	now := time.Now()
	var added int
	for _, ad := range addrs {
		// 0 port is invalid
		if ad.Port == 0 {
			continue
		}
		// Discard own client
		if ad.IP.IsLoopback() && ad.Port == d.listenPort {
			continue
		} else if d.clientIP.Equal(ad.IP) {
			continue
		}
		if externalip.IsExternal(ad.IP) {
			continue
		}
		if d.blocklist != nil && d.blocklist.Blocked(ad.IP) {
			continue
		}
		p := &peerAddr{
			addr:      ad,
			timestamp: now,
			source:    source,
			priority:  peerpriority.Calculate(ad, d.clientAddr()),
		}
		//Priority of peer
		//Calculate priority of peers for peer to peer connection
		//Different peer priority calculation criteria can be there (geographic, hardware properties of the client, ISP )
		//Peer priority criteria which is most optimum for p2p network can be determined using static rules, Machine Learning Algorithm
		item := d.peerByPriority.ReplaceOrInsert(p)
		if item != nil {
			prev := item.(*peerAddr)
			d.peerByTime[prev.index] = p
			p.index = prev.index
			d.countBySource[prev.source]--
		} else {
			d.peerByTime = append(d.peerByTime, p)
			p.index = len(d.peerByTime) - 1
		}
		added++
	}
	d.filterNils()
	sort.Sort(byTimestamp(d.peerByTime))
	d.countBySource[source] += added


	delta := d.peerByPriority.Len() - d.maxItems
	//Peer priority list has exceeded maximum capacity and it needs to be trimmed
	if delta > 0 {
		
		d.removeExcessItems(delta)
		d.filterNils()
		d.countBySource[source] -= delta
	}
	if len(d.peerByTime) != d.peerByPriority.Len() {
		panic("addr list data structures not in sync")
	}
}

func (d *AddrList) filterNils() {
	b := d.peerByTime[:0]
	for _, x := range d.peerByTime {
		if x != nil {
			b = append(b, x)
			x.index = len(b) - 1
		}
	}
	d.peerByTime = b
}

//Remove excess items 
//Delete the farthest node (calculated according to peer priority algorithm) and trim till maximum capacity
func (d *AddrList) removeExcessItems(delta int) {
	for i := 0; i < delta; i++ {
		d.peerByPriority.Delete(d.peerByTime[i])
		d.peerByTime[i] = nil
	}
}

func (d *AddrList) clientAddr() *net.TCPAddr {
	ip := *d.clientIP
	if ip == nil {
		ip = net.IPv4(0, 0, 0, 0)
	}
	return &net.TCPAddr{
		IP:   ip,
		Port: d.listenPort,
	}
}
