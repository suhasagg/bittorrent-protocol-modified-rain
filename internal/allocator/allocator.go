package allocator

import (
	"github.com/cenkalti/rain/internal/metainfo"
	"github.com/cenkalti/rain/internal/storage"
)

// Allocator allocates files on the disk.
// Just for file serialization - File disk map Algorithm
type Allocator struct {
	Files       []File
	HasExisting bool
	HasMissing  bool
	Error       error

	closeC chan struct{}
	doneC  chan struct{}
}

// File on the disk.
type File struct {
	Storage storage.File
	Name    string
}

// Progress about the allocation.
type Progress struct {
	AllocatedSize int64
}

// New returns a new Allocator.
func New() *Allocator {
	return &Allocator{
		closeC: make(chan struct{}),
		doneC:  make(chan struct{}),
	}
}

// Close the Allocator.
func (a *Allocator) Close() {
	close(a.closeC)
	<-a.doneC
}

// Run the Allocator.
// Per torrent meta information (how many file pieces are there for a torrent for a file)
// On the basis of this, how much storage is needed for a torrent is computed for a peer
func (a *Allocator) Run(info *metainfo.Info, sto storage.Storage, progressC chan Progress, resultC chan *Allocator) {
	defer close(a.doneC)

	defer func() {
		if a.Error != nil {
			for _, f := range a.Files {  
				if f.Storage != nil {
					f.Storage.Close()
				}
			}
		}
		select {
		case resultC <- a:
		case <-a.closeC:
		}
	}()

	var allocatedSize int64
	a.Files = make([]File, len(info.Files))
	for i, f := range info.Files {
		var sf storage.File
		var exists bool
		// 
		sf, exists, a.Error = sto.Open(f.Path, f.Length)
		if a.Error != nil {
			return
		}
		a.Files[i] = File{Storage: sf, Name: f.Path}
		if exists {
			a.HasExisting = true
		} else {
			a.HasMissing = true
		}
		//Cumulative storage needed by a torrent 
		allocatedSize += f.Length
		a.sendProgress(progressC, allocatedSize)
	}
}

func (a *Allocator) sendProgress(progressC chan Progress, size int64) {
	select {
    //Progress channel is used to track how much torrent size has been written to disk currently (how much torrent has been downloaded)
	case progressC <- Progress{AllocatedSize: size}:
	case <-a.closeC:
		return
	}
}
