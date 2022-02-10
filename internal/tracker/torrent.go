package tracker

// Torrent contains fields that are sent in an announce request.
type Torrent struct {
	BytesUploaded   int64
	BytesDownloaded int64
	BytesLeft       int64
	InfoHash        [20]byte //(Torrent for which file) //Fetching file parts from the peers
	PeerID          [20]byte // Peers specific
	Port            int
}
