package timelord

import (
	"fmt"
	"sync"
	"time"

	"github.com/bouk/monkey"
)

var mut = &sync.Mutex{}
var isTimeLord bool

// TimeLord is capable of manipulating time by monkey-patching time.Now.
type TimeLord struct {
	guard  *monkey.PatchGuard
	offset time.Duration
}

// New returns a new TimeLord with a time offset of 0. It is not safe to create
// more than one TimeLord at a time, so subsequent calls to New will return an
// error.
func New() (*TimeLord, error) {
	mut.Lock()
	defer mut.Unlock()
	if isTimeLord {
		return nil, fmt.Errorf("timelord.New has already been called (only one TimeLord can be created at a time)")
	}
	tl := &TimeLord{}
	tl.Warp()
	return tl, nil
}

// Warp monkey-patches time.Now to return the current time + offset.
func (tl *TimeLord) Warp() {
	mut.Lock()
	defer mut.Unlock()
	if tl.guard != nil {
		tl.guard.Restore()
	}
	tl.guard = monkey.Patch(time.Now, func() time.Time {
		tl.guard.Unpatch()
		defer tl.guard.Restore()

		return time.Now().Add(tl.offset)
	})
}

// Unwarp restores time.Now to its original functionality.
func (tl *TimeLord) Unwarp() {
	tl.guard.Unpatch()
}

// SetOffset sets the time offset. You do not need to call Warp again after
// setting the offset.
func (tl *TimeLord) SetOffset(d time.Duration) {
	tl.offset = d
}
