package timelord

import (
	"fmt"
	"sync"
	"time"

	"github.com/bouk/monkey"
)

var globalMut = &sync.Mutex{}
var timeLordExists bool

// TimeLord is capable of manipulating time by monkey-patching time.Now.
type TimeLord struct {
	guard  *monkey.PatchGuard
	offset time.Duration
}

// New returns a new TimeLord with a time offset of 0. It is not safe to create
// more than one TimeLord at a time, so subsequent calls to New will return an
// error, unless this TimeLord is destroyed.
func New() (*TimeLord, error) {
	globalMut.Lock()
	defer globalMut.Unlock()
	if timeLordExists {
		return nil, fmt.Errorf("timelord.New has already been called (only one TimeLord can be created at a time)")
	}
	tl := &TimeLord{}
	tl.warp()
	timeLordExists = true
	return tl, nil
}

// SetOffset sets the time offset, which has the effect of making time.Now
// return the current time + offset.
func (tl *TimeLord) SetOffset(d time.Duration) {
	tl.offset = d
}

// Destroy destroys the current TimeLord, restores time.Now to its default
// behavior, and allows a new one to be created via New.
func (tl *TimeLord) Destroy() {
	globalMut.Lock()
	defer globalMut.Unlock()
	tl.unwarp()
	tl.guard = nil
	timeLordExists = false
}

// warp monkey-patches time.Now to return the current time + offset.
func (tl *TimeLord) warp() {
	tl.guard = monkey.Patch(time.Now, func() time.Time {
		tl.guard.Unpatch()
		defer tl.guard.Restore()

		return time.Now().Add(tl.offset)
	})
}

// unwarp restores time.Now to its original functionality.
func (tl *TimeLord) unwarp() {
	if tl.guard != nil {
		tl.guard.Unpatch()
	}
}
